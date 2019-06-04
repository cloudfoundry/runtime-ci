#!/usr/bin/env ruby

require 'hashdiff'
require_relative './binary_changes.rb'
require_relative './renderer.rb'

updates = BinaryUpdates.new(
  'cf-deployment-concourse-tasks-latest-release/dockerfiles/cf-deployment-concourse-tasks/Dockerfile',
  'cf-deployment-concourse-tasks/dockerfiles/cf-deployment-concourse-tasks/Dockerfile'
)

puts Renderer.new.render(
  binary_updates: updates,
  task_updates: TaskUpdates.generate(
    'cf-deployment-concourse-tasks-latest-release',
    'cf-deployment-concourse-tasks'
  )
)
