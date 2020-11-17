require 'tmpdir'
require 'rspec'
require 'fileutils'
require_relative './release_changes.rb'
require 'webmock/rspec'

describe 'ReleaseUpdates' do
  before(:all) do
    @current_work_dir = Dir.pwd
    @tmp_work_dir = Dir.mktmpdir('cf-deployment-test')

    Dir.chdir(@tmp_work_dir)
  end

  after(:all) do
    Dir.chdir(@current_work_dir)
    FileUtils.rm_rf(@tmp_work_dir) if File.exists?(@tmp_work_dir)
  end

  let(:ru) { ReleaseUpdates.new }

  describe 'load_from_files' do
    before(:all) do
      FileUtils.mkdir_p('cf-deployment-main/operations')
      FileUtils.mkdir_p('cf-deployment-release-candidate/operations')
    end

    subject(:updates) do
      ReleaseUpdates.load_from_files(filename, opsfile: opsfile)
    end

    before do
      stub_request(:get, /github.com/)

      File.open(File.join('cf-deployment-main', filename), 'w') do |f|
        f.write(file_contents_main)
      end

      File.open(File.join('cf-deployment-release-candidate', filename), 'w') do |f|
        f.write(file_contents_rc)
      end
    end

    let(:filename) { 'cf-deployment.yml' }
    let(:opsfile) { false }
    let(:file_contents_main) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.1.0
  url: https://bosh.io/d/github.com/org/release-1?v=1.1.0
- name: release-2
  version: 2.1.0
  url: https://bosh.io/d/github.com/org/release-2?v=2.1.0
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
  url: https://bosh.io/d/github.com/org/release-1?v=1.2.0
- name: release-2
  version: 2.2.0
  url: https://bosh.io/d/github.com/org/release-2?v=2.2.0
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

      expect(release_1_update.old_url).to eq 'https://github.com/org/release-1/releases/tag/v1.1.0'
      expect(release_1_update.new_url).to eq 'https://github.com/org/release-1/releases/tag/v1.2.0'

      expect(release_2_update.old_url).to eq 'https://github.com/org/release-2/releases/tag/v2.1.0'
      expect(release_2_update.new_url).to eq 'https://github.com/org/release-2/releases/tag/v2.2.0'
    end

    it 'reads the inputs, and returns a list of stemcell updates' do
      stemcell_update = updates.get_update_by_name('ubuntu-trusty')
      expect(stemcell_update.old_version).to eq 1
      expect(stemcell_update.new_version).to eq 2
    end

    context('when the file is an ops-file') do
      let(:filename) { 'operations/ops-file.yml' }
      let(:opsfile) { true }
      let(:file_contents_main) do
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
    url: https://bosh.io/d/github.com/cloudfoundry-incubator/hwc-buildpack-release?v=2.3.4
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
    url: https://bosh.io/d/github.com/cloudfoundry-incubator/hwc-buildpack-release?v=2.4.0
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

      context 'when the ops-file replaces releases and stemcells rather than appending them' do
        let(:file_contents_main) do
<<-HEREDOC
- type: replace
  path: /releases/name=garden-windows?
  value:
    name: garden-windows
    version: 0.6.0
- type: replace
  path: /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/cc/stacks?
- type: replace
  path: /releases/name=hwc-buildpack?
  value:
    name: hwc-buildpack
    version: 2.3.4
- type: replace
  path: /stemcells/name=windows2012R?
  value:
    name: windows2012R2
    version: 1
HEREDOC
        end

        let(:file_contents_rc) do
<<-HEREDOC
- type: replace
  path: /releases/name=garden-windows?
  value:
    name: garden-windows
    version: 0.7.0
- type: replace
  path: /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/cc/stacks?
- type: replace
  path: /releases/name=hwc-buildpack?
  value:
    name: hwc-buildpack
    version: 2.4.0
- type: replace
  path: /stemcells/name=windows2012R
  value:
    name: windows2012R2
    version: 2

HEREDOC
        end

        it 'returns the release and stemcell updates of the replacements binaries' do
          release_1_update = updates.get_update_by_name('garden-windows')
          release_2_update = updates.get_update_by_name('hwc-buildpack')
          stemcell_update = updates.get_update_by_name('windows2012R2')

          expect(release_1_update.old_version).to eq '0.6.0'
          expect(release_1_update.new_version).to eq '0.7.0'

          expect(release_2_update.old_version).to eq '2.3.4'
          expect(release_2_update.new_version).to eq '2.4.0'

          expect(stemcell_update.old_version).to eq 1
          expect(stemcell_update.new_version).to eq 2
        end
      end
    end

    context 'when the old version of the file is empty' do
      let(:file_contents_main) { "" }
      context 'and the new version is not empty' do
        let(:file_contents_rc) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.2.0
stemcells:
- os: ubuntu-trusty
  version: 2
HEREDOC
        end

        it 'views the newly-introduced releases as additive updates' do
          stemcell_update = updates.get_update_by_name('ubuntu-trusty')
          expect(stemcell_update.old_version).to eq nil
          expect(stemcell_update.new_version).to eq 2

          release_update = updates.get_update_by_name('release-1')
          expect(release_update.old_version).to eq nil
          expect(release_update.new_version).to eq "1.2.0"
        end
      end

      context 'and the new version is empty' do
        let(:file_contents_rc) { "" }

        it 'includes no information about the release' do
          expect(updates.get_update_by_name('ubuntu-trusty')).to be_nil
          expect(updates.get_update_by_name('release-1')).to be_nil
        end
      end
    end

    context 'when the new version of the file is empty' do
      let(:file_contents_rc) { "" }
      context 'and the old version is not empty' do
        let(:file_contents_main) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.2.0
stemcells:
- os: ubuntu-trusty
  version: 2
HEREDOC
        end

        it 'views the removed releases as negative updates' do
          stemcell_update = updates.get_update_by_name('ubuntu-trusty')
          expect(stemcell_update.old_version).to eq 2
          expect(stemcell_update.new_version).to eq nil

          release_update = updates.get_update_by_name('release-1')
          expect(release_update.old_version).to eq "1.2.0"
          expect(release_update.new_version).to eq nil
        end
      end

      context 'and the new version is empty' do
        let(:file_contents_main) { "" }

        it 'includes no information about the release' do
          expect(updates.get_update_by_name('ubuntu-trusty')).to be_nil
          expect(updates.get_update_by_name('release-1')).to be_nil
        end
      end
    end

    context 'when the old version of the file does not exist' do
      before do
        File.delete(File.join('cf-deployment-main', filename))

        File.open(File.join('cf-deployment-release-candidate', filename), 'w') do |f|
          f.write(file_contents_rc)
        end
      end

      it 'treats the release as having a new version' do
          stemcell_update = updates.get_update_by_name('ubuntu-trusty')
          expect(stemcell_update.old_version).to eq nil
          expect(stemcell_update.new_version).to eq 2

          release_update = updates.get_update_by_name('release-1')
          expect(release_update.old_version).to eq nil
          expect(release_update.new_version).to eq "1.2.0"
      end
    end

    context 'when the new version of the file does not exist' do
      before do
        File.delete(File.join('cf-deployment-release-candidate', filename))

        File.open(File.join('cf-deployment-main', filename), 'w') do |f|
          f.write(file_contents_main)
        end
      end

      it 'treats the release as having been deleted' do

          stemcell_update = updates.get_update_by_name('ubuntu-trusty')
          expect(stemcell_update.old_version).to eq 1
          expect(stemcell_update.new_version).to eq nil

          release_update = updates.get_update_by_name('release-1')
          expect(release_update.old_version).to eq "1.1.0"
          expect(release_update.new_version).to eq nil
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

  describe '#convert_bosh_io_to_github_url' do
    before do
      File.open(File.join('cf-deployment-main', filename), 'w') do |f|
        f.write(file_contents_main)
      end
    end

    let(:filename) { 'cf-deployment.yml' }
    let(:opsfile) { false }
    let(:file_contents_main) do
<<-HEREDOC
releases:
- name: release-1
  version: 1.1.0
  url: https://bosh.io/d/github.com/org/release-1?v=1.1.0
stemcells:
- os: ubuntu-trusty
  version: 1
HEREDOC
    end

    context 'release located outside of bosh.io' do
      let(:url) { 'https://not-bosh-at-all.io/some-release' }
      it 'should return an empty url' do
        expect(ru.send(:convert_bosh_io_to_github_url, url)).to eq(nil)
      end
    end
  end

  describe '#find_correct_github_url_with_prefix' do

    context 'handles github redirect' do
      before do
        stub_request(:get, 'https://github.com/org/release-1/releases/tag/v1.1.0')
          .to_return(status: 301, headers: {'location': 'https://github.com/new-org/release-1/releases/tag/v1.1.0'})
      end

      it 'same tag prefix' do
        stub_request(:get, 'https://github.com/new-org/release-1/releases/tag/v1.1.0')
          .to_return(status: 200)
        release_update = ReleaseUpdates.load_from_files('cf-deployment.yml')
        release = release_update.get_update_by_name("release-1")

        expect(release.old_url).to eq 'https://github.com/new-org/release-1/releases/tag/v1.1.0'
      end

      it 'different tag prefix' do
        stub_request(:get, 'https://github.com/new-org/release-1/releases/tag/v1.1.0')
          .to_return(status: 404)
        stub_request(:get, 'https://github.com/new-org/release-1/releases/tag/1.1.0')
          .to_return(status: 200)

        release_update = ReleaseUpdates.load_from_files('cf-deployment.yml')
        release = release_update.get_update_by_name("release-1")

        expect(release.old_url).to eq 'https://github.com/new-org/release-1/releases/tag/1.1.0'
      end
    end
  end

  describe '#generate_github_uri_from_bosh_io' do
    let(:input_uri) { URI('https://bosh.io/d/github.com/org/release-1?v=1.1.0') }
    let(:output_uri) { URI('https://github.com/org/release-1/releases/tag/1.1.0') }

    it 'convert bosh.io url to github.com' do
        expect(ru.send(:generate_github_uri_from_bosh_io, input_uri)).to eq(output_uri)
    end
  end

  describe '#inject_tag_prefix_into_github_uri' do
    let(:input_uri) { URI('https://github.com/org/release-1/releases/tag/1.1.0') }
    let(:output_uri) { URI('https://github.com/org/release-1/releases/tag/1.1.0') }
    let(:output_uri_with_tag_prefix) { URI('https://github.com/org/release-1/releases/tag/v1.1.0') }

    it 'inject empty prefix' do
        expect(ru.send(:inject_tag_prefix_into_github_uri, input_uri, '')).to eq(output_uri)
    end

    it 'inject v prefix' do
        expect(ru.send(:inject_tag_prefix_into_github_uri, input_uri, 'v')).to eq(output_uri_with_tag_prefix)
    end
  end

  describe '#merge!' do
    it 'merges two sets of release updates' do
      updates1 = ReleaseUpdates.new
      updates2 = ReleaseUpdates.new

      updates1.load_change(['-', '', {"name"=>"release", "version" => 1}])
      updates1.load_change(['+', '', {"name"=>"release", "version" => 2}])

      updates2.load_change(['-', '', {"name" => "stemcell", "version" => 1}])

      updates1.merge!(updates2)

      release_update = updates1.get_update_by_name("release")
      expect(release_update.old_version).to eq 1
      expect(release_update.new_version).to eq 2

      stemcell_update = updates1.get_update_by_name("stemcell")
      expect(stemcell_update.old_version).to eq 1
      expect(stemcell_update.new_version).to be_nil
    end
  end
end
