require 'yaml'

raise "Usage: ruby template_the_manifest.rb MANIFEST_YML MANIFEST_CONFIG_YML" unless ARGV.size == 2


def traverse(node, &blk)
  case node
  when Array
    node.each do |e|
      traverse(e, &blk)
    end

  when Hash
    node.each do |_k,v|
      traverse(v, &blk)
    end
  else
    blk.call(node)
  end
  node
end

manifest_body = YAML.load_file(ARGV[0])
manifest_replacements = YAML.load_file(ARGV[1])
manifest_replacements.each do |replace_key, replace_value|
  traverse(manifest_body) do |value|
    value.gsub!(replace_key, replace_value) if value.is_a?(String)
  end
end

puts YAML.dump(manifest_body)
