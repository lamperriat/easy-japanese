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


# curl -X POST "http://localhost:8080/api/words/book_1/submit" \
# -H "Content-Type: application/json" \
# -H "X-API-KEY: TEST_USE_API_KEY" \
# -d '{
#     "kanji": "test",
#     "chinese": "test",
#     "example": []
# }'


# curl -X POST "http://localhost:8080/api/words/book_1/delete" \
# -H "Content-Type: application/json" \
# -H "X-API-KEY: TEST_USE_API_KEY" \
# -d '{
#     "id": 5,
#     "kanji": "test",
#     "chinese": "test",
#     "example": []
# }'

echo add
curl -X POST "http://localhost:8080/api/reading-material/add" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "id": 0,
    "content": "test",
    "chinese": "test cn"
}'

echo edit
curl -X POST "http://localhost:8080/api/reading-material/edit" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "id": 1,
    "content": "test edited",
    "chinese": "test cn edited"
}'

echo "search fail"
curl -X GET "http://localhost:8080/api/reading-material/search?page=1&query=wrong" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"

echo "search success"
curl -X GET "http://localhost:8080/api/reading-material/search?page=1&query=test" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"

echo "get"
curl -X GET "http://localhost:8080/api/reading-material/get" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"

echo "delete"
curl -X POST "http://localhost:8080/api/reading-material/delete" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "id": 1,
    "content": "test edited",
    "chinese": "test cn edited"
}'

echo "get"
curl -X GET "http://localhost:8080/api/reading-material/get" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"

curl -X GET "http://localhost:8080/api/words/book_1/get" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"

curl -X POST "http://localhost:8080/api/words/all/accurate-search" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY" \
-d '{
    "chinese": "string",
    "example": [
      {
        "chinese": "string",
        "example": "string"
      }
    ],
    "hiragana": "string",
    "id": 0,
    "kanji": "中国人",
    "katakana": "string",
    "type": "string"
}'

curl -X GET "http://localhost:8080/api/words/all/fuzzy-search?query=ちゅ" \
-H "Content-Type: application/json" \
-H "X-API-KEY: TEST_USE_API_KEY"
