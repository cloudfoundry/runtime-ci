require 'rspec'
require 'fileutils'
require_relative './release_changes.rb'

describe 'ReleaseUpdates' do
  describe 'load_from_files' do
    before(:all) do
      FileUtils.mkdir_p('cf-deployment-master/operations')
      FileUtils.mkdir_p('cf-deployment-release-candidate/operations')
    end

    subject(:updates) do
      ReleaseUpdates.load_from_files(filename, opsfile: opsfile)
    end

    before do
      File.open(File.join('cf-deployment-master', filename), 'w') do |f|
        f.write(file_contents_master)
      end

      File.open(File.join('cf-deployment-release-candidate', filename), 'w') do |f|
        f.write(file_contents_rc)
      end
    end

    let(:filename) { 'cf-deployment.yml' }
    let(:opsfile) { false }
    let(:file_contents_master) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.1.0
- name: release-2
  version: 2.1.0
stemcells:
- os: ubuntu-trusty
  version: 1
HEREDOC
    end

    let(:file_contents_rc) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.2.0
- name: release-2
  version: 2.2.0
stemcells:
- os: ubuntu-trusty
  version: 2
HEREDOC
    end

    it 'reads the given file in the two inputs, and returns the release updates' do
      release_1_update = updates.get_update_by_name('release-1')
      release_2_update = updates.get_update_by_name('release-2')

      expect(release_1_update.old_version).to eq '1.1.0'
      expect(release_1_update.new_version).to eq '1.2.0'

      expect(release_2_update.old_version).to eq '2.1.0'
      expect(release_2_update.new_version).to eq '2.2.0'
    end

    it 'reads the inputs, and returns a list of stemcell updates' do
      stemcell_update = updates.get_update_by_name('ubuntu-trusty')
      expect(stemcell_update.old_version).to eq 1
      expect(stemcell_update.new_version).to eq 2
    end

    context('when the file is an ops-file') do
      let(:filename) { 'operations/ops-file.yml' }
      let(:opsfile) { true }
      let(:file_contents_master) do
<<-HEREDOC
- type: replace
  path: /releases/-
  value:
    name: garden-windows
    version: 0.6.0
- type: replace
  path: /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/cc/stacks?
- type: replace
  path: /releases/-
  value:
    name: hwc-buildpack
    version: 2.3.4
- type: replace
  path: /stemcells/-
  value:
    name: windows2012R2
    version: 1
HEREDOC
      end

      let(:file_contents_rc) do
<<-HEREDOC
- type: replace
  path: /releases/-
  value:
    name: garden-windows
    version: 0.7.0
- type: replace
  path: /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/cc/stacks?
- type: replace
  path: /releases/-
  value:
    name: hwc-buildpack
    version: 2.4.0
- type: replace
  path: /stemcells/-
  value:
    name: windows2012R2
    version: 2

HEREDOC
      end

      it 'reads the given file in the two inputs, and returns the release updates' do
        release_1_update = updates.get_update_by_name('garden-windows')
        release_2_update = updates.get_update_by_name('hwc-buildpack')

        expect(release_1_update.old_version).to eq '0.6.0'
        expect(release_1_update.new_version).to eq '0.7.0'

        expect(release_2_update.old_version).to eq '2.3.4'
        expect(release_2_update.new_version).to eq '2.4.0'
      end

      it 'reads the inputs, and returns a list of stemcell updates' do
        stemcell_update = updates.get_update_by_name('windows2012R2')
        expect(stemcell_update.old_version).to eq 1
        expect(stemcell_update.new_version).to eq 2
      end
    end
  end

  describe '#load_change' do
    subject(:updates) { ReleaseUpdates.new }

    let(:version) { rand(10).to_s }
    let(:name) { 'capi-release' }
    let(:change) do
      [
        op,
        '[20]',
        {
          'name' => name,
          'version' => version
        }
      ]
    end

    context 'when the operation is "+"' do
      let(:op) { '+' }
      it 'saves the version as new_version' do
        updates.load_change(change)
        expect(updates.get_update_by_name(name).new_version).to eq(version)
      end
    end

    context 'when the operation is "-"' do
      let (:op) { '-' }
      it 'saves the version as old_version' do
        updates.load_change(change)
        expect(updates.get_update_by_name(name).old_version).to eq version
      end
    end

    context 'when a second change for the same release occurs' do
      let(:change1) do
        [
          '-',
          '[20]',
          {
            'name' => name,
            'version' => '26'
          }
        ]
      end

      let(:change2) do
        [
          '+',
          '[20]',
          {
            'name' => name,
            'version' => '27'
          }
        ]
      end

      it 'saves the old and new versions together' do
        subject.load_change(change1)
        subject.load_change(change2)
        expect(subject.get_update_by_name(name).new_version).to eq '27'
        expect(subject.get_update_by_name(name).old_version).to eq '26'
      end
    end
  end
end
