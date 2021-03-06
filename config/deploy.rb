require "pathname"
require "whenever/capistrano"

# config valid only for current version of Capistrano
lock '3.4.0'

set :application, 'my_app_name'
set :repo_url, 'https://github.com/YoshitsuguFujii/slackbot'

# Default branch is :master
# ask :branch, `git rev-parse --abbrev-ref HEAD`.chomp
set :branch, ENV['BRANCH'] || "master"

# Default deploy_to directory is /var/www/my_app_name
# set :deploy_to, '/var/www/my_app_name'

# Default value for :scm is :git
# set :scm, :git

# Default value for :format is :pretty
# set :format, :pretty

# Default value for :log_level is :debug
# set :log_level, :debug

# Default value for :pty is false
# set :pty, true

# Default value for :linked_files is []
# set :linked_files, fetch(:linked_files, []).push('config/database.yml', 'config/secrets.yml')
set :linked_files, fetch(:linked_files, []).push('tmp.pid')

# Default value for linked_dirs is []
# set :linked_dirs, fetch(:linked_dirs, []).push('log', 'tmp/pids', 'tmp/cache', 'tmp/sockets', 'vendor/bundle', 'public/system')
set :linked_dirs, fetch(:linked_dirs, []).push('log', 'todo_list/stored_files')

# Default value for default_env is {}
# set :default_env, { path: "/opt/ruby/bin:$PATH" }

# Default value for keep_releases is 5
# set :keep_releases, 5

set :rbenv_type, :system
set :rbenv_ruby, '2.2.3'

set :whenever_identifier, ->{ "#{fetch(:application)}_#{fetch(:stage)}" }
set :whenever_roles, :slack # 対象ロール
set :whenever_environment, :production # 対象環境

namespace :deploy do
  after :restart, :clear_cache do
    on roles(:web), in: :groups, limit: 3, wait: 10 do
      # Here we can do anything such as:
      # within release_path do
      #   execute :rake, 'cache:clear'
      # end
    end
  end

  desc 'kill'
  task :kill do
    on roles(:slack) do |host|
      #if test "[ ! -d #{current_path}/tmp.pid ]"
      if test "[ -f #{shared_path}/tmp.pid ]"
        execute "cd #{shared_path} && kill -9 `cat tmp.pid` > /dev/null "" 2>&1"
        execute "cd #{shared_path} && rm tmp.pid"
      end
    end
  end
  before 'deploy:starting', 'deploy:kill'

  desc 'create_shared_directory'
  task :create_shared_directory do
    on roles(:slack) do |host|
      unless test "[ -d #{shared_path}/log ]"
        execute :mkdir, '-pv', "#{shared_path}/log"
      end

      unless test "[ -d #{shared_path}/todo_list/stored_files ]"
        execute :mkdir, '-pv', "#{shared_path}/todo_list/stored_files"
      end
    end
  end

  desc 'create_shared_files'
  task :create_shared_files do
    on roles(:slack) do |host|
      unless test "[ -f #{shared_path}/tmp.pid ]"
        execute "touch #{shared_path}/tmp.pid"
      end
    end
  end

  before 'deploy:starting', 'deploy:create_shared_directory'
  before 'deploy:starting', 'deploy:create_shared_files'

  desc 'gom install'
  task :gom_install do
    on roles(:slack) do |host|
      execute "cd #{release_path} && gom install"
    end
  end
  before 'deploy:published', 'deploy:gom_install'

  desc 'upload files'
  task :upload_files do
    on roles(:slack) do |host|
      upload!('slackbot_responder/word.yml', "#{release_path}/slackbot_responder/word.yml")
      upload!('twitterbot/watch_user.yml', "#{release_path}/twitterbot/watch_user.yml")
      upload!('twitterbot/watch_word.yml', "#{release_path}/twitterbot/watch_word.yml")

      Pathname.glob("api-mock/*.yml").each do |file|
        upload!(file.to_s, "#{release_path}/api-mock")
      end
    end
  end
  before 'deploy:published', 'deploy:upload_files'

  desc 'build'
  task :build do
    on roles(:slack) do |host|
      execute "cd #{release_path} && GOOS=linux GOARCH=amd64 gom build *.go"
    end
  end
  before 'deploy:published', 'deploy:build'

  desc 'run'
  task :run do
    on roles(:slack) do |host|
      execute "cd #{current_path} && bin/start.sh"
    end
  end
  before 'deploy:finished', 'deploy:run'
end
