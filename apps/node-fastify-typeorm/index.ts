import f from 'fastify'
import { Column, Entity, PrimaryGeneratedColumn, DataSource } from "typeorm"

const fastify = f()

@Entity({ name: 'users' })
class User {
  @PrimaryGeneratedColumn()
  id!: number

  @Column({ type: 'varchar' })
  name!: string

  @Column({ type: 'int' })
  age!: number
}

const dataSource = new DataSource({
  type: 'postgres',
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  username: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  synchronize: false,
  logging: false,
  entities: [User],
  extra: { max: 100 }
})

dataSource.initialize()

fastify.get('/', async (request, reply) => {
  const userRepository = dataSource.getRepository(User)

  const user = new User()
  user.name = 'Test!'
  user.age = 1
  await userRepository.save(user)

  const userDb = await userRepository.findOneBy({ id: user.id })

  userDb!.age += 1
  await userRepository.save(userDb!)

  await userRepository.delete(userDb!)

  return userDb
})

const start = async () => {
  try {
    await fastify.listen({ port: 3000 }, (err, address) => {
      if (err) throw err
      console.log(`Listening ${address}`)
    })
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start()

