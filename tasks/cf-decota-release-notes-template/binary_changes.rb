require 'hashdiff'

class BinaryUpdate
  attr_accessor :old_version, :new_version
end

class BinaryUpdates
  class << self
    def load_from_file(filename)
      binary_current_list = collect_binaries("cf-deployment-concourse-tasks/#{filename}")
      binary_latest_release_list = collect_binaries("cf-deployment-concourse-tasks-latest-release/#{filename}")

      binary_updates = BinaryUpdates.new
      change_set = HashDiff.diff(binary_latest_release_list, binary_current_list)
      change_set.each do |change|
        binary_updates.load_change(change)
      end

      binary_updates
    end

    private

    def collect_binaries(dockerfile_path)
      binaries = `grep 'ENV .*_version' #{dockerfile_path} | awk 'gsub("_version", "")' | awk '{print $2 ":" $3}'`
      binaries.split("\n").map do |b|
        {
          "name" => b.split(":")[0],
          "version" => b.split(":")[1]
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

    name = change[2]['name']
    version = change[2]['version']

    binary_update = @updates[name] || BinaryUpdate.new

    if op == '+'
      binary_update.new_version = version
    elsif op == '-'
      binary_update.old_version = version
    end

    @updates[name] = binary_update
  end

  def get_update_by_name(binary_name)
    @updates[binary_name]
  end

  def each
    @updates.each do |binary_name, binary_update|
      yield binary_name, binary_update
    end
  end

  def merge!(updates2)
    updates2.each do |binary_name, binary_update|
      @updates[binary_name] = binary_update
    end
  end
end
