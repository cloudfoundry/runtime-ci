require 'rspec'
require 'fileutils'
require_relative './release_changes.rb'

describe 'ReleaseUpdates' do
  describe 'load_from_files' do
    before(:all) do
      FileUtils.mkdir_p('cf-deployment-master')
      FileUtils.mkdir_p('cf-deployment-release-candidate')
    end

    let(:filename) { 'cf-deployment.yml' }
    before do
      File.open(File.join('cf-deployment-master', filename), 'w') do |f|
        f.write(
<<-HEREDOC
releases:
- name: release-1
  version: 1.1.0
- name: release-2
  version: 2.1.0
HEREDOC
        )
      end

      File.open(File.join('cf-deployment-release-candidate', filename), 'w') do |f|
        f.write(
<<-HEREDOC
releases:
- name: release-1
  version: 1.2.0
- name: release-2
  version: 2.2.0
HEREDOC
        )
      end
    end

    it 'reads the given file in the two inputs, and returns the release updates' do
      updates = ReleaseUpdates.load_from_files(filename)
      release_1_update = updates.get_update_by_name('release-1')
      release_2_update = updates.get_update_by_name('release-2')

      expect(release_1_update.old_version).to eq '1.1.0'
      expect(release_1_update.new_version).to eq '1.2.0'

      expect(release_2_update.old_version).to eq '2.1.0'
      expect(release_2_update.new_version).to eq '2.2.0'
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
