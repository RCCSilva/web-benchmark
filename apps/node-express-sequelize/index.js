'use strict';

const express = require('express')
const { Sequelize, DataTypes } = require('sequelize');

const app = express()
const port = 3000

const url = `postgres://${process.env.DB_USER}:${process.env.DB_PASSWORD}@${process.env.DB_HOST}:${process.env.DB_PORT}/${process.env.DB_NAME}`
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

app.get('/', async (req, res) => {
  try {
    const user = await User.create({ name: 'Test!', age: 1 })
    const userDb = await User.findByPk(user.id)
    userDb.age += 1
    await userDb.save()
    await userDb.destroy()
    res.status(200).send(userDb.toJSON())
  } catch (err) {
    console.error(err)
    res.status(500).send({})
  }
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
