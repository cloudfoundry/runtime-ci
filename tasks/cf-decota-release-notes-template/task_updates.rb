require 'digest'

class TaskUpdates
  def generate(latest_release_dir, master_dir)
    latest_release_task_list = TaskFinder.find_tasks(latest_release_dir)
    master_task_list = TaskFinder.find_tasks(master_dir)

    new_task_list = master_task_list - latest_release_task_list
    deleted_task_list = latest_release_task_list - master_task_list
    common_task_list = latest_release_task_list & master_task_list
    updated_task_list = []

    common_task_list.each do |task|
     latest_release_task_shasum = task_shasum(latest_release_dir, task)
     master_task_shasum = task_shasum(master_dir, task)
      if latest_release_task_shasum != master_task_shasum
        updated_task_list << task
      end
    end

    {
      "new" => new_task_list,
      "updated" => updated_task_list,
      "deleted" => deleted_task_list,
    }
  end

  private

  def task_shasum(dir, task)
    Digest::SHA1.hexdigest(
      Dir.glob(
        File.join(dir, task, '*')
      ).map do |file|
        Digest::SHA1.hexdigest(read_file(file))
      end.join
    )
  end

  def read_file(file)
    if !file.match(/.*task.yml/)
      return File.read(file)
    end

    File.readlines(file).reject do |line|
      line.match(/.*tag: v\d+.\d+.\d+/)
    end.join
  end
end
