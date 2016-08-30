#!/usr/bin/env ruby
require 'yaml'

input_filename = './deployment-manifest/' + ENV.fetch('DEPLOYMENT_MANIFEST')
output_filename = './modified-deployment-manifest/' + ENV.fetch('MODIFIED_DEPLOYMENT_MANIFEST')
properties = YAML.load_file(input_filename)
release_name = ENV.fetch('RELEASE_NAME')
release_to_modify = properties['releases'].find do |release|
  release['name'] == release_name
end
properties['releases'].delete(release_to_modify)
properties['releases'] << {
    'name' => release_name,
    'version' => 'latest'
}

File.open(output_filename, 'w') do |file|
  YAML.dump(properties, file)
end
