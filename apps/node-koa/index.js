'use strict';

const Koa = require('koa');
const app = new Koa();

app.use(async ctx => {
  ctx.body = { status: 'ok' };
  ctx.status = 200
});

app.listen(3000);
