#!/usr/bin/env ruby
require 'yaml'
require 'active_support'

input_filename = 'bosh-lite-stub/' + ENV.fetch('BOSH_LITE_STUB_PATH')
output_filename = 'extended-bosh-lite-stub/' + ENV.fetch('EXTENDED_BOSH_LITE_STUB_PATH')

properties = YAML.load_file(input_filename)
properties.deep_merge('property_override' => {'nsync' => {'diego_privileged_containers' => true}})

File.open(output_filename, 'w') do |file|
  YAML.dump(properties, file)
end
