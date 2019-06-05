#!/usr/bin/env ruby

require 'hashdiff'
require_relative './binary_changes.rb'
require_relative './task_updates.rb'
require_relative './renderer.rb'

puts Renderer.new.render(
  binary_updates: BinaryUpdates.new(
  'cf-deployment-concourse-tasks-latest-release/dockerfiles/cf-deployment-concourse-tasks/Dockerfile',
  'cf-deployment-concourse-tasks/dockerfiles/cf-deployment-concourse-tasks/Dockerfile'
  ),
  task_updates: TaskUpdates.new(
    'cf-deployment-concourse-tasks-latest-release',
    'cf-deployment-concourse-tasks'
  )
)
