#!/usr/bin/env ruby
require 'socket'      # Sockets are in standard library

s = TCPSocket.open( 'localhost', 4000)

100.times do |i|
  s.puts 'GET name'
  puts s.gets
end
s.close               # Close the socket when done
