1.  使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

   #### 命令

   ```
   redis-benchmark -d 10 -t get,set
   redis-benchmark -d 20 -t get,set
   redis-benchmark -d 50 -t get,set
   redis-benchmark -d 100 -t get,set
   redis-benchmark -d 200 -t get,set
   redis-benchmark -d 1000 -t get,set
   redis-benchmark -d 5000 -t get,set
   ```

   #### 结果

   | value大小 | 10万requests  SET | 10万requests GET |
   | --------- | ----------------- | ---------------- |
   | 10        | 28344.67 r/s      | 29120.56 r/s     |
   | 20        | 29256.88          | 29895.37         |
   | 50        | 30039.05          | 29325.51         |
   | 100       | 27078.26          | 27948.57         |
   | 200       | 29403.12          | 30321.41         |
   | 1000      | 30339.81          | 29274.01         |
   | 5000      | 25588.54          | 26680.90         |

   #### 结论

   value在1000以内，测试性能没有什么影响，在打到5000的大小后，性能突降



2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

   

未写入情况

```
used_memory:691336
used_memory_human:675.13K
```

1w 10字节

```
used_memory:2393360 - 691336
used_memory_human:2.28M
平均 (2393360 - 691336)/10000 = 170.2024
```

50w  10字节

```
used_memory:81079944 - 691336
used_memory_human:77.32M
平均 (81079944 - 691336)/500000 = 160.777216
```

1w 1000字节

```
used_memory:12313360 - 691336
used_memory_human:11.74M
平均 (12313360 - 691336) / 10000 = 1162.2024
```

50w的key   1000字节

```
used_memory:577080064 - 691336
used_memory_human:550.35M
(577080064 - 691336) / 500000 = 1152.777456
```

1w的key  5000的字节

```
used_memory:83993360 -691336 
used_memory_human:80.10M
(83993360 -691336) / 10000 = 8330.2024
```



50w的key 5000字节

机器限制，实际上只测了225036

```
used_memory:1873642224 - 691336
used_memory_human:1.74G
(1873642224 - 691336) / 225036 = 8322.89450577
```

 从以上结果来看

（对于1w和50w的存储）存储key的多少对平均每个key的影响不大

对于字节，从10字节到1000字节，再到5000字节，从10字节到1000字节，数据量增长了100倍，但是存储增长并不明显，而对于5000字节，数据量增长5倍，增储量增长不止5倍，说明大的value对存储占用较高