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
  'routing',
  'cf',
]

deployment_configuration_path = ENV.fetch('DEPLOYMENT_CONFIGURATION_PATH')
deployment_manifest_path = ENV.fetch("DEPLOYMENT_MANIFEST_PATH")
stemcell_name = ENV.fetch("STEMCELL_NAME")

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

stemcell_url = File.read("stemcell/url")
stemcell_version = File.read("stemcell/version")
stemcell_sha1 = File.read("stemcell/sha1")

stemcell_metadata = {
  'stemcells' => [
    {
      'name' => stemcell_name,
      'url' => stemcell_url,
      'version' => stemcell_version,
      'sha1' => stemcell_sha1
    }
  ]
}

releases = YAML.dump(releases_metadata).gsub("---\n", '')
stemcells = YAML.dump(stemcell_metadata).gsub("---\n", '')

deployment_configuration = File.read("deployment-configuration/#{deployment_configuration_path}")
updated_deployment_manifest = "#{deployment_configuration}\n#{releases}\n#{stemcells}"

File.open("deployment-manifest/#{deployment_manifest_path}", 'w') do |f|
  f.write(updated_deployment_manifest)
end
