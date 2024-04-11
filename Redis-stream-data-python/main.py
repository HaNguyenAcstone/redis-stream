import logging
from flask import Flask, request
from redis import Redis
import json
import time

app = Flask(__name__)
redis_conn = Redis(host='192.168.38.128', port=6379, db=0)

# Thiết lập logging level cho Flask app
app.logger.setLevel(logging.INFO)

# Thiết lập format cho log messages
formatter = logging.Formatter('[%(asctime)s] %(levelname)s in %(module)s: %(message)s')

# Tạo một file handler để ghi log vào console
file_handler = logging.FileHandler('app.log')
file_handler.setLevel(logging.INFO)
file_handler.setFormatter(formatter)
app.logger.addHandler(file_handler)

# Tên của Redis Stream
stream_name = 'Redis_Streams_AcstOne'

# Tạo Redis Stream (nếu chưa tồn tại)
if not redis_conn.exists(stream_name):
    # Tạo Redis Stream
    redis_conn.xadd(stream_name, {'init': 'start'})

# Push Data 1 Time
def push_orders_to_redis(orders):
    for order in orders:
        # Convert đơn hàng thành JSON trước khi đẩy vào Redis Stream
        order_json = json.dumps(order)
        # Đẩy đơn hàng vào Redis Stream với key là 'order'
        redis_conn.xadd(stream_name, {'order': order_json})

# Route để nhận yêu cầu push đơn hàng
@app.route('/push_orders', methods=['GET'])
def push_orders():
    # Nhận tham số từ URL
    key = request.args.get('key')
    value = request.args.get('value')

    # Kiểm tra nếu key và value không tồn tại
    if not key or not value:
        return 'Invalid parameters', 400

    # Lấy giá trị value và chuyển đổi sang kiểu int
    num_orders = int(value)

    # Tạo danh sách đơn hàng
    orders = [{'order_id': i} for i in range(1, num_orders + 1)]

    # Thời điểm bắt đầu xử lý yêu cầu
    start_time = time.time()

    # Gửi đơn hàng vào Redis Stream
    push_orders_to_redis(orders)

    # Thời gian hoàn thành yêu cầu
    end_time = time.time()
    completion_time = end_time - start_time

    # Ghi log về thời gian hoàn thành của yêu cầu
    app.logger.info(f"Request completed in {completion_time} seconds")

    return f'Orders pushed to Redis Stream successfully. Completion time: {completion_time} seconds'

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')