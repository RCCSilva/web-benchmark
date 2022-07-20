'use strict';

const Koa = require('koa');
const { Sequelize, DataTypes } = require('sequelize');
const app = new Koa();

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


app.use(async ctx => {
  try {
    const user = await User.create({ name: 'Test!', age: 1 })
    const userDb = await User.findByPk(user.id)
    userDb.age += 1
    await userDb.save()
    await userDb.destroy()

    ctx.body = userDb.toJSON()
  } catch (err) {
    console.error(err)
    ctx.status = 500
    ctx.body = {}
  }
});

app.listen(3000);
