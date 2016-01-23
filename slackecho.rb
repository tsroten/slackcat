class Slackecho < Formula
  desc "Simple command-line utility to post messages to Slack."
  homepage "https://github.com/tsroten/slackecho"
  url "https://github.com/tsroten/slackecho/archive/v1.1.tar.gz"
  version "1.1"
  sha256 "9705e3fde7dad407c88aec9cd6c888b5087d1a6d9a8d975bef8b6781ed939d9c"

  depends_on "go"
  depends_on "slackcat"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/slackecho-1.1-#{platform}-amd64" => "slackecho"

    puts "Slackecho is installed. If you don't already have Slackcat,\n"
    puts "install and configure Slackcat with:"
    puts "  brew install slackcat"
  end

  test do
    assert_equal(0, "/usr/local/bin/slackecho")
  end
end
