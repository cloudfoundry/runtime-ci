require 'rspec'
require_relative './variable_updates.rb'

describe 'VariableUpdates' do
  describe '.load_from_files' do
    before(:all) do
      FileUtils.mkdir_p('cf-deployment-master/operations')
      FileUtils.mkdir_p('cf-deployment-release-candidate/operations')
    end

    subject(:updates) do
      VariableUpdates.load_from_files(filename)
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
    let(:file_contents_master) do
<<-HEREDOC
variables:
- name: gets-removed
  type: password
- name: remains
  type: ssh
- name: gets-updated
  type: ssh
- name: ca-cert
  type: certificate
  options:
    is_ca: true
    common_name: common-name
- name: leaf-cert
  type: certificate
  options:
    ca: ca_cert
    common_name: leaf-common-name
    extended_key_usage:
    - server-auth
    alternative_names:
    - alternative-name

HEREDOC
    end

    let(:file_contents_rc) do
<<-HEREDOC
variables:
- name: remains
  type: ssh
- name: gets-updated
  type: password
- name: gets-added
  type: rsa
- name: ca-cert
  type: certificate
  options:
    is_ca: true
    common_name: common-name
- name: leaf-cert
  type: certificate
  options:
    ca: ca_cert
    common_name: updated-leaf-common-name
    extended_key_usage:
    - server-auth
    - new-ext-key-usage
    alternative_names:
    - alternative-name

HEREDOC
    end

    it 'reads the given file in the two inputs, and returns the variable updates' do
      gets_removed_update = updates.get_update_by_name('gets-removed')
      expect(gets_removed_update.state).to eq :removed

      remains_update = updates.get_update_by_name('remains')
      expect(remains_update).to be_nil

      gets_updated_update = updates.get_update_by_name('gets-updated')
      expect(gets_updated_update.state).to eq :updated

      gets_added_update = updates.get_update_by_name('gets-added')
      expect(gets_added_update.state).to eq :added

      ca_cert_update = updates.get_update_by_name('ca-cert')
      expect(ca_cert_update).to be_nil

      leaf_cert_update = updates.get_update_by_name('leaf-cert')
      expect(leaf_cert_update.state).to eq :updated
    end
  end

  describe '#load_change' do
    subject(:updates) { VariableUpdates.new }
    context 'when the operation is `+`' do
      let(:new_change) do
        [
          '+',
          '[0]',
          {'name' => 'variable-name', 'type' => 'password'}
        ]
      end

      context 'and the operation previously did not exist' do
        it 'adds an update with state :added' do
          updates.load_change(new_change)

          update = updates.get_update_by_name('variable-name')
          expect(update.state).to eq :added
          expect(update.type).to eq 'password'
        end
      end

      context 'and the previous operation was `-`' do
        before do
          updates.load_change([
            '-',
            '[0]',
            {'name' => 'variable-name', 'type' => 'ssh'}
          ])
        end

        it 'updates the object with state :updated and updates the type' do
          updates.load_change(new_change)

          update = updates.get_update_by_name('variable-name')
          expect(update.state).to eq :updated
          expect(update.type).to eq 'password'
        end
      end

      context 'and the previous operation was `+`' do
        before do
          updates.load_change([
            '+',
            '[0]',
            {'name' => 'variable-name', 'type' => 'ssh'}
          ])
        end

        it 'raises' do
          expect { updates.load_change(new_change) }.to raise_error 'Disallowed No-op'
        end
      end
    end

    context 'when the operation is `-`' do
      let(:new_change) do
        [
          '-',
          '[0]',
          {'name' => 'variable-name', 'type' => 'password'}
        ]
      end

      context 'and the operation previously did not exist' do
        it 'adds an update with state :removed' do
          updates.load_change(new_change)

          update = updates.get_update_by_name('variable-name')
          expect(update.state).to eq :removed
          expect(update.type).to eq 'password'
        end
      end

      context 'and the previous operation was `-`' do
        before do
          updates.load_change([
            '-',
            '[0]',
            {'name' => 'variable-name', 'type' => 'ssh'}
          ])
        end

        it 'raise' do
          expect { updates.load_change(new_change) }.to raise_error 'Disallowed No-op'
        end
      end

      context 'and the previous operation was `+`' do
        before do
          updates.load_change([
            '+',
            '[0]',
            {'name' => 'variable-name', 'type' => 'ssh'}
          ])
        end

        it 'updates the object with state :updated and updates the type' do
          updates.load_change(new_change)

          update = updates.get_update_by_name('variable-name')
          expect(update.state).to eq :updated
          expect(update.type).to eq 'password'
        end
      end

    end
  end
end
