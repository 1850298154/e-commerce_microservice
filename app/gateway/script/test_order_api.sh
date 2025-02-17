#!/bin/bash

# 设置基础 URL
BASE_URL="http://localhost:8080"

# 定义测试函数
run_test() {
    local test_num=$1
    echo "开始运行测试 #${test_num}"
    
    # 1. 创建订单
    echo "Testing Place Order API (#${test_num})..."
    ORDER_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8080/orders" \
      -H "Content-Type: application/json" \
      -d '{
        "user_currency": "USD",
        "address": {
            "street_address": "123 Main St",
            "city": "Boston",
            "state": "MA",
            "country": "USA",
            "zip_code": 12345
        },
        "email": "test@example.com",
        "order_items": [
            {
                "product_id": 1,
                "quantity": 2,
                "cost": 9.99
            }
        ]
    }')

    echo "Place Order Response (#${test_num}): $ORDER_RESPONSE"

    # 从响应中提取订单ID
    ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.order')

    # 2. 获取订单列表
    echo "\nTesting List Orders API (#${test_num})..."
    curl -s -X GET "${BASE_URL}/orders"

    # 3. 标记订单已支付
    echo "\nTesting Mark Order as Paid API (#${test_num})..."
    curl -s -X PUT "${BASE_URL}/orders/${ORDER_ID}/paid" \
      -H "Content-Type: application/json" \
      -d '{
    }'

    # 4. 更新订单
    echo "\nTesting Update Order API (#${test_num})..."
    curl -s -X PUT "${BASE_URL}/orders/${ORDER_ID}" \
      -H "Content-Type: application/json" \
      -d '{
        "new_address": {
            "street_address": "456 New St",
            "city": "New York",
            "state": "NY",
            "country": "USA",
            "zip_code": 54321
        },
        "new_email": "newemail@example.com",
        "new_order_items": [
            {
                "product_id": 1,
                "quantity": 3,
                "cost": 9.99
            }
        ]
    }'

    # 5. 取消订单
    echo "\nTesting Cancel Order API (#${test_num})..."
    curl -s -X DELETE "${BASE_URL}/orders/${ORDER_ID}" \
      -H "Content-Type: application/json" \
      -d '{
        "timed_cancel": false,
        "cancel_time": 0
    }'
    
    echo "完成测试 #${test_num}"
}

# 并发运行10个测试实例
for i in {1..10}; do
    run_test $i &
done

# 等待所有后台任务完成
wait

echo "所有测试完成"