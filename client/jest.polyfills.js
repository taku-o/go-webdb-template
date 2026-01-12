/**
 * Additional polyfills for Jest + MSW v2
 * jest-fixed-jsdom handles most Fetch API globals
 */

const { TextDecoder, TextEncoder } = require('node:util')

Object.defineProperties(globalThis, {
  TextDecoder: { value: TextDecoder },
  TextEncoder: { value: TextEncoder },
})

// BroadcastChannel is required by msw v2
const { BroadcastChannel } = require('node:worker_threads')

Object.defineProperties(globalThis, {
  BroadcastChannel: { value: BroadcastChannel },
})
