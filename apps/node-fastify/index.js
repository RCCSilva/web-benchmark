'use strict';

// Require the framework and instantiate it
const fastify = require('fastify')({ logger: false })

// Declare a route
fastify.get('/', async (request, reply) => {
  return { status: 'ok' }
})

// Run the server!
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
