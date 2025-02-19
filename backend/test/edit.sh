# curl -X POST "http://localhost:8080/api/words/book_1/edit" \
# -H "Content-Type: application/json" \
# -H "X-API-KEY: TEST_USE_API_KEY" \
# -d '{
#     "id": 1,
#     "kanji": "中国人",
#     "chinese": "中国人",
#     "example": [
#         {"example": "私は中国人です", "chinese": "我是中国人"}
#     ]
# }'

# curl -X PUT "http://localhost:8080/api/words/book_1/edit" \
# -H "Content-Type: application/json" \
# -H "X-API-KEY: your_secret_key" \
# -d '{"id": 99999}'

# curl -X PUT "http://localhost:8080/api/words/wrong_dict/edit" \
# -H "Content-Type: application/json" \
# -H "X-API-KEY: your_secret_key" \
# -d '{"id": 1}'


curl -X POST "http://localhost:8080/api/words/book_1/submit" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "kanji": "test",
    "chinese": "test",
    "example": []
}'


curl -X POST "http://localhost:8080/api/words/book_1/delete" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "id": 5,
    "kanji": "test",
    "chinese": "test",
    "example": []
}'