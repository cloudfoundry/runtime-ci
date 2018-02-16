require 'rspec'
require_relative './ops_file_finder.rb'

describe 'OpsFileFinder' do
  REPO_DIR = 'test-cf-deployment'
  before(:all) do
      FileUtils.mkdir_p(REPO_DIR)
  end

  context 'when there are ops-files in the lop-level operations directory' do
    before do
      File.open("#{REPO_DIR}/operations/ops.yml", "w") {}
      File.open("#{REPO_DIR}/operations/ops2.yml", "w") {}
      File.open("#{REPO_DIR}/operations/README.md", "w") {}
    end

    it 'returns the file names of the ops-files without any additional path' do
      ops = OpsFileFinder.find_ops_files(REPO_DIR)

      expect(ops).to include 'ops.yml'
      expect(ops).to include 'ops2.yml'
    end

    it 'does not return any files that are not yaml files' do
      ops = OpsFileFinder.find_ops_files(REPO_DIR)

      expect(ops).not_to include 'README.md'
    end
  end

  context 'when there is another directory in the operations directory' do
    before do
      File.open("#{REPO_DIR}/operations/experimental/ops.yml", "w") {}
      File.open("#{REPO_DIR}/operations/experimental/ops2.yml", "w") {}
    end

    it 'does not include the directory in the list of files' do
      ops = OpsFileFinder.find_ops_files(REPO_DIR)

      expect(ops).not_to include 'experimental'
    end

  end
end
