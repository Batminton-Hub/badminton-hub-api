# Project freetime Badminton-Hub

# -- Goal --
- เพื่อให้ง่ายต่อการ setup 
- ติดตามผลได้จริง 
- matainace ได้ง่าย

# -- Setup .env --
> MODE = โหมดการทำงาน หากเป็น deverlop function ที่มีการ random จะเป็น default
> SERVER_PORT = port ที่ใช้รันระบบ
> DB_NAME = ชื่อ Database หลัก
> MONGO_DB_URL = url ที่ใช้ต่อกับ mongoDB
> KEY_BEARER_TOKEN = key ที่ใช้เข้า-ถอดรหัส BearerToken
> KEY_HASH_AUTH = key ที่ใช้เข้า hash ที่ติดไปกับ BearToken
> KEY_HASH_MEMBER = key ที่ใช้สร้าง hash ประจำตัว member
> KEY_HASH_PASSWORD = key ที่ใช้เข้ารหัส password ของ member ให้เป็นความลับ
> GOOGLE_CALLBACK_LOGIN_URL = url callback ของระบบ login ผ่าน google
> GOOGLE_CALLBACK_REGISTER_URL = url callback ของระบบ register ผ่าน google
> GOOGLE_CLIENT_ID = clinet ID ที่ออกโดย google cloud console
> GOOGLE_CLIENT_SECRET = clinet secret ที่ออกโดย google cloud console
> REDIS_CACHE_ADDR = url เชื่อมต่อกับ redis
> REDIS_CACHE_PASSWORD = password ของ redis
> REDIS_CACHE_DB = db หลักของ redis
> DEFAULT_AES_IV = ค่า defalut ของ iv ที่ใช้กับการเข้ารหัสแบบ AES
> DEFAULT_GOOGLE_STATE = ค่า defalut ของ state ที่ใช้ร่วมกับ OAuth v2
> BEARER_TOKEN_EXP = lifetime ของ BearerToken
> TRACER_SERVER_URL = url สำหรับ server tracer



