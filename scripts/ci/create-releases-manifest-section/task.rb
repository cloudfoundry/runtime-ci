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

deployment_configuration_path = ENV.fetch('DEPLOYMENT_CONFIGURATION_PATH')
deployment_manifest_path = ENV.fetch("DEPLOYMENT_MANIFEST_PATH")

releases_metadata = release_names.map do |release_name|
  release_resource = "#{release_name}-release"

  url = File.read("#{release_resource}/url").strip
  version = File.read("#{release_resource}/version").strip
  sha1 = File.read("#{release_resource}/sha1").strip

  {
    'name' => release_name,
    'url' => url,
    'version' => version,
    'sha1' => sha1
  }
end

puts "Updated releases"
releases = YAML.dump(releases_metadata).gsub("---\n", '')

deployment_configuration = File.read("deployment-configuration/#{deployment_configuration_path}")
updated_deployment_manifest = "#{deployment_configuration}\n#{releases}"

File.open("deployment-manifest/#{deployment_manifest_path}", 'w') do |f|
  f.write(updated_deployment_manifest)
end
