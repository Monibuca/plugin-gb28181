# PS结构介绍

### Key Point 
1. 一帧原始数据只能包含一个`ps`, `map`,`sys`(关键帧需要)等结构.当一帧数据超过pes最大的负载时(`0xFFFF`)，则一帧数据可能有多个`pes`结构 
2. rtp结构是包含在封装的最外层，所以 一帧数据或者一片pes结构，可能会被分成多个rtp包
3. 各个结构中可能都有扩充字段，在解析的时候可以通过扩充字段来解析到表示的长度，从而得到真实的原始数据的长度
| ps | system header | psm | pesv | pesa |
`pesv`为视频pes包， `pesa`为音频pes包
**非关键帧的ps包结构一般不包含系统头和PSM**

针对PSM中流类型的取值如下：
| 流类型     |   十六进制值|
|-----------|-----------|
| MPEG-4    | 0x10      |
| H.264     | 0x1B      |
| SVAC      | 0x80      |
| G.711     | 0x90      |
| G.722.1   | 0x92      |
| G.723.1   | 0x93      |
| G.729     | 0x99      |
| SVAC 音频  | 0x9B      |
具体可以参看`ISO/IEC 13818-1:2000`

### 各个结构的头字段的说明

#### PS 
| 字段名 | 比特位长度 | 含义|
|------------------|:-----:|----------------------------------------|
| pack_start_code  | 32 |  起始码字段`0x000001BA` 标志为一包的开始 |
| marker_bit |2|标记位字段2位字段，取值`01`| 
| system_clock_reference_base (32..30) |3|系统时钟参考字段|
| marker_bit | 1| 标记位字段取值`1` |
| system_clock_reference_base(29..15) |15|系统时钟参考字段 |
| marker_bit |1 |标记位字段取值`1` |
| system_clock_reference_base(14..0) |15|系统时钟参考字段 |
| marker_bit |1 |标记位字段取值`1`|	
| system_clock_reference_extension |9|系统始终参考字段 |
| marker_bit |1|标记位字段取值`1`|
| program_mux_rate |22|节目复合速率字段 |
| marker_bit |1| 标记記位字段取值`1`|
| marker_bit |1| 标记位字段取值`1`|
| reserved  |5| 填充字段 |
| pack_stuffing_length |3 |包填充长度字段|



#### System 
| 字段名 | 比特位长度 | 含义|
|------------------|:-----:|----------------------------------------|
| system_header_start_code | 32| 开始码 `0x000001BB` |
| header_length |16|該字段後的系統標題的字節長度 |
| marker_bit |1| 标记位字段取值`1` |
| rate_bound |22|速率界限字段 |
| marker_bit |1|标记位字段取值`1`| 
| audio_bound |6| 音频界限字段 |
| fixed_flag |1|固定标志字段 |
| CSPS_flag |1|CSPS标志字段 |
| system_audio_lock_flag |1|系統音频锁定标志字段 |
| system_video_lock_flag |1| 系統视频锁定标志字段|
| marker_bit |1| 标记位字段取值`1`| 
| vedio_bound |5| 视频界限字段 |
| packet_rate_restriction_flag|1 |分组速率限制标记字段|
| reserved_bits |7| 保留位字段|
| stream_id |8| 流标志字段  一般`0xe0`指视频，`0xc0`音频|
| marker_bit |2|取值`11` |
| P-STD_buffer_bound_scale |1b| P-STD界限比例字段|
|P-STD_buffer_size_bound |13|P-STD大小界限字段|
注：可能会有多个stream , 通过stream_id的第一个bit判断是否存在stream

#### Map
| 字段名 | 比特位长度 | 含义|
|------------------|:-----:|----------------------------------------| 
| packet start code |32|映射流起始码标志字段|
| program stream map length |16|节目流映射长度字段 |
| current_next_indicator |1|当前下一個指示符字段 |
| reserved |2|填充字段 |
| program_stream_map_version|5|节目流映射版本字段 |
| reserved |7|预留填充字段|
| marker_bit |1|标记位字段取值`1`|
| program_stream_info_length |16| 节目流信息长度字段 |
| elementary_stream_map_length|16|基本流映射长度字段 |
| stream_type |8| 流类型字段 `0x1b H264`， `0x24 H265`|
| elementary_stream_id |8|视频取值`0xe0-0xef`，通常为`0xe0`,音频取值`0xc0-0xdf`，通常为`0xc0`|
| elementary_stream_info_length |16|基本流信息长度字段 |
| CRC_32 |32| CRC字段|


#### PES 
| 字段名 | 比特位长度 | 含义|
|------------------|:-----:|----------------------------------------| 
| packet_start_code_prefix |24| 分组起码码`0x000001` |
| stream_id  |8|流标志字段stream_id，其中0x(C0~DF)指音频，0x(E0~EF)视频|
| PES_packet_length |16| PES分组长度字段|
| marker_bit|2| 标志位，`10`|
| PES_scrambling_control |2|PES加扰控制字段 |
| PES_priority |1| PES优先级字段 |
| data_alignment_indicator |1|数据对其指示符字段|
| copyright |1| 盘权字段|
| original_or_copy |1| 原始或拷贝字段 |
| PTS_DTS_flags |2| `PTS DTS`标志字段|
| ESCR_flag |1|ESCR标志字段| 
| ES_rate_flag |1|ES速率标志字段| 
| DSM_trick_mode_flag |1|DSM特技方式标志字段 |
| additional_copy_info_flag |1|附加版权信息标志字段 |
| PES_CRC_flag |1|PES CRC标志字段 |
| PES_extension_flag|1| PES扩张标志字段 |
| PES_header_data_length| 8| PES标题数据长度字段 |

注：后面会根据此长度来填写 pts， dts等
