package edge_tts

import "errors"

// Errors for the Edge TTS project.

// UnknownResponse error for unknown server responses
var UnknownResponse = errors.New("unknown response received from the server")

// UnexpectedResponse error for unexpected server responses
//
// This hasn't happened yet, but it's possible that the server will change its response format in the future.
var UnexpectedResponse = errors.New("unexpected response received from the server")

// NoAudioReceived error when no audio is received from the server
var NoAudioReceived = errors.New("no audio received from the server")

// WebSocketError error for WebSocket errors
var WebSocketError = errors.New("WebSocket error occurred")
