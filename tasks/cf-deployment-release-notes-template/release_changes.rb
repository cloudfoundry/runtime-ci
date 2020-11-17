require 'yaml'
require 'hashdiff'
require 'net/http'
require 'uri'

class ReleaseUpdate
  attr_accessor :old_version, :new_version, :old_url, :new_url
end

class ReleaseUpdates
  class << self
    def load_from_files(filename, opsfile: false)
      release_candidate = load_yaml_file('cf-deployment-release-candidate', filename)
      release_candidate_releases_list = collect_releases_and_stemcells(release_candidate, opsfile: opsfile)

      main_branch = load_yaml_file('cf-deployment-main', filename)
      main_releases_list = collect_releases_and_stemcells(main_branch, opsfile: opsfile)

      release_updates = ReleaseUpdates.new
      change_set = HashDiff.diff(main_releases_list, release_candidate_releases_list)
      change_set.each do |change|
        release_updates.load_change(change)
      end

      release_updates
    end

    private

    def load_yaml_file(input_name, filename)
      filepath = File.join(input_name, filename)

      return nil unless File.exist? filepath

      file_text = File.read(filepath)
      YAML.load(file_text) || nil
    end

    def collect_releases_and_stemcells(manifest, opsfile: false)
      return [] if manifest.nil? || manifest.empty?
      return manifest['releases'] + manifest['stemcells'] unless opsfile

      manifest.select do |op|
        op['type'] == 'replace' && (op['path'].start_with?('/releases') || op['path'].start_with?('/stemcells'))
      end.collect do |op|
        {
          "name" => op["value"]["name"] || op['value']['os'],
          "version" => op["value"]["version"],
          "url" => op["value"]["url"]
        }
      end
    end
  end

  def initialize
    @updates = {}
  end

  def load_change(change)
    op = change[0]
    if op == "~"
      return
    end

    name = change[2]['name'] || change[2]['os']
    version = change[2]['version']
    url = change[2]['url'] ? convert_bosh_io_to_github_url(change[2]['url']) : nil

    release_update = @updates[name] || ReleaseUpdate.new

    if op == '+'
      release_update.new_version = version
      release_update.new_url = url
    elsif op == '-'
      release_update.old_version = version
      release_update.old_url = url
    end

    @updates[name] = release_update
  end

  def convert_bosh_io_to_github_url(bosh_url)
    bosh_uri = URI(bosh_url)

    return nil unless bosh_uri.host.match 'bosh.io'

    generated_github_uri = generate_github_uri_from_bosh_io(bosh_uri)

    find_correct_github_url_with_prefix(generated_github_uri)
  end

  def find_correct_github_url_with_prefix(generated_github_uri)
    tag_prefixes = ['v', '']

    tag_prefixes.each do |prefix|
      prefixed_github_uri = inject_tag_prefix_into_github_uri(generated_github_uri, prefix)
      return nil unless prefixed_github_uri

      generated_url_response = Net::HTTP.get_response(prefixed_github_uri)

      ok = generated_url_response.code == '200'
      redirect = generated_url_response.code == '301'

      return prefixed_github_uri.to_s if ok

      generated_url_response.header['location']
      new_location = generated_url_response.header['location']

      return find_correct_github_url_with_prefix(URI(new_location)) if redirect
    end

    nil
  end

  def generate_github_uri_from_bosh_io(bosh_uri)
    github_string = bosh_uri.path.sub('/d/', '')
    host, *path = github_string.split('/')
    version = URI.decode_www_form(bosh_uri.query).assoc('v').last

    project_path = []
    project_path.concat(path)
    project_path << 'releases'
    project_path << 'tag'
    project_path << version

    URI::HTTPS.build(host: host,
                     path: '/' + project_path.join('/'))
  end

  def inject_tag_prefix_into_github_uri(uri, prefix)
    raise 'github.com url expected' unless uri&.host&.match 'github.com'

    new_uri = uri.clone
    new_uri.path = uri.path.gsub(%r{tag/v?}, 'tag/' + prefix)
    new_uri
  end

  def get_update_by_name(release_name)
    @updates[release_name]
  end

  def each
    @updates.each do |release_name, release_update|
      yield release_name, release_update
    end
  end

  def merge!(updates2)
    updates2.each do |release_name, release_update|
      @updates[release_name] = release_update
    end
  end
end
