'use strict';

const fastify = require('fastify')({ logger: false })
const { Sequelize, DataTypes } = require('sequelize');


const url =
  `postgres://${process.env.DB_USER}:${process.env.DB_PASSWORD}@${process.env.DB_HOST}:${process.env.DB_PORT}/${process.env.DB_NAME}`

const sequelize = new Sequelize(url, { pool: { max: 100 } })

console.log(`Connected to database ${url}`)

const User = sequelize.define('User', {
  id: {
    type: DataTypes.BIGINT,
    primaryKey: true,
    autoIncrement: true
  },
  name: DataTypes.STRING,
  age: DataTypes.INTEGER,
}, {
  timestamps: false,
  tableName: 'users'
});

fastify.get('/', async (request, reply) => {

  const user = await User.create({ name: 'Test!', age: 1 })
  const userDb = await User.findByPk(user.id)
  userDb.age += 1
  await userDb.save()
  await userDb.destroy()

  return userDb.toJSON()
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

