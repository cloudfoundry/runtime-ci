#!/usr/bin/ruby -w

require 'json'

ips = ARGV[0]
destination_file = ARGV[1]

asgs=[]

ips.split(/\s/).each do |ip|
  asgs << {protocol: 'tcp', destination: ip}
end

File.open(destination_file,"w") do |f|
  f.write(asgs.to_json)
end
