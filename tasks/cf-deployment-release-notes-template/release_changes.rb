require 'yaml'
require 'hashdiff'

class ReleaseUpdate
  attr_accessor :old_version, :new_version
end

class ReleaseUpdates
  class << self
    def load_from_files(filename, opsfile: false)
      release_candidate_text = File.read(File.join('cf-deployment-release-candidate', filename))
      release_candidate = YAML.load(release_candidate_text)

      master_text = File.read(File.join('cf-deployment-master', filename))
      master = YAML.load(master_text)

      if opsfile
        master_releases_list = filter_release_changes(master)
        release_candidate_releases_list = filter_release_changes(release_candidate)

        master_stemcells_list = filter_stemcell_changes(master)
        release_candidate_stemcells_list = filter_stemcell_changes(release_candidate)
      else
        master_releases_list = master['releases']
        release_candidate_releases_list = release_candidate['releases']

        master_stemcells_list = master['stemcells']
        release_candidate_stemcells_list = release_candidate['stemcells']
      end

      release_updates = ReleaseUpdates.new
      changeSet = HashDiff.diff(master_releases_list, release_candidate_releases_list)
      changeSet.each do |change|
        release_updates.load_change(change)
      end

      changeSet = HashDiff.diff(master_stemcells_list, release_candidate_stemcells_list)
      changeSet.each do |change|
        release_updates.load_change(change)
      end

      release_updates
    end

    private

    def filter_release_changes(ops_list)
      ops_list.select do |op|
        op['type'] == 'replace' && op['path'] == '/releases/-'
      end.collect do |op|
        {
          "name" => op["value"]["name"],
          "version" => op["value"]["version"]
        }
      end
    end

    def filter_stemcell_changes(ops_list)
      ops_list.select do |op|
        op['type'] == 'replace' && op['path'] == '/stemcells/-'
      end.collect do |op|
        {
          "name" => op["value"]["name"],
          "version" => op["value"]["version"]
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

    release_update = @updates[name] || ReleaseUpdate.new

    if op == '+'
      release_update.new_version = version
    elsif op == '-'
      release_update.old_version = version
    end

    @updates[name] = release_update
  end

  def get_update_by_name(release_name)
    @updates[release_name]
  end

  def each
    @updates.each do |release_name, release_update|
      yield release_name, release_update
    end
  end
end
