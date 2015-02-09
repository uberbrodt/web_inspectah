require 'typhoeus'

hydra = Typhoeus::Hydra.hydra
good_requests = 0
bad_requests = 0

File.open('request_data.txt').each do |request|
    t_request = Typhoeus::Request.new(request)
    t_request.on_complete do |response|
      if response.code == 204
        good_requests = good_requests+1
      elsif response.code == 403
        puts response.request.url
        bad_requests = bad_requests+1
      end
    end
    hydra.queue t_request
end

hydra.run
puts "good requests: #{good_requests}"
puts "bad requests: #{bad_requests}"
