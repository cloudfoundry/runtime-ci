#!/usr/bin/env ruby

require 'hashdiff'
require_relative './binary_changes.rb'
require_relative './renderer.rb'

updates = BinaryUpdates.load_from_file('dockerfiles/cf-deployment-concourse-tasks/Dockerfile')
puts Renderer.new.render(binary_updates: updates)
