import express from 'express'
import { Column, Entity, PrimaryGeneratedColumn, DataSource } from "typeorm"

const app = express()
const port = 3000

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

app.get('/', async (req, res) => {
  const userRepository = dataSource.getRepository(User)

  const user = new User()
  user.name = 'Test!'
  user.age = 1
  await userRepository.save(user)

  const userDb = await userRepository.findOneBy({ id: user.id })

  userDb!.age += 1
  await userRepository.save(userDb!)

  await userRepository.delete(userDb!)

  res.status(200).send(userDb)
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
