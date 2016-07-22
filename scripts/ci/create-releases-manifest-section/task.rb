#!/usr/bin/env ruby

require 'yaml'

release_names=[
  'consul',
  'diego',
  'etcd',
  'loggregator',
  'cf-mysql',
  'uaa',
  'garden-linux',
  'cflinuxfs2-rootfs',
  'cf',
]

releases_metadata = release_names.map do |release_name|
  release_resource = "#{release_name}-release"

  url = File.read("#{release_resource}/url")
  version = File.read("#{release_resource}/version")
  sha1 = File.read("#{release_resource}/sha1")

  {
    'name' => release_name,
    'url' => url,
    'version' => version,
    'sha1' => sha1
  }
end

puts YAML.dump(releases_metadata)
