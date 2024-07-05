# int_version = 0x17000841        # 7.0.8.65
# int_version = 0x17000C2B        # 7.0.12.43
int_version = 0x1800322B        # 8.0.50.43
# 转16进制字符串
hex_version = hex(int_version)

# 去掉0x
hex_str = hex_version[2:]

# 把第一个字符（最高位）替换为 0
new_hex_str = "0" + hex_str[1:]

# 转回10进制
new_hex_num = int(new_hex_str, 16)

# 按位还原版本号
major = (new_hex_num >> 24) & 0xFF
minor = (new_hex_num >> 16) & 0xFF
patch = (new_hex_num >> 8) & 0xFF
build = (new_hex_num >> 0) & 0xFF

# 拼接版本号
wx_version = "{}.{}.{}.{}".format(major, minor, patch, build)
print(wx_version)