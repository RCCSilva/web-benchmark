'use strict';

const express = require('express')

const app = express()
const port = 3000


app.get('/', async (req, res) => {
  res.status(200).send({ status: 'ok' })
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
