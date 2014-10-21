create table if not exists petLog_%d.t_login_%s (date  datetime,uin  int unsigned,ret  int,version  int,ip  varchar(255),login_type  int,level  int,plat_id  int,vip  varchar(255),index(date),index(version, ret, uin),index(date, ret, uin),index(vip, ret, uin),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_lazywriter_online_%s (date  datetime,segment  int,zoneid  int,client_onlie  int,all_online  int) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_zonesvr_online_%s (date  datetime,zoneid  int,segment  int,online  int) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_zone_login_%s (date  datetime,uin  int unsigned,ret  int,zoneid  int,areaid  int,new_player  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_pay_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,vip  int,payid  int,bill_id  varchar(255),paychannel  int,pay_uin  int unsigned,pay_charge  int,level  int,plat_id  int,index(payid),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_buy_%s (date  datetime,uin  int unsigned,petid  varchar(255),bill_id  varchar(255),ybqb  varchar(255),vip  int,charge  int,get_uin  int,get_petid  varchar(255),src_id  int,cwq_charge  int,plat_id  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_buy_detail_%s (date  datetime,uin  int unsigned,petid  varchar(255),bill_id  varchar(255),ybqb  varchar(255),commodity_id  int,num  int,price  int,src_id  int,vip  int,plat_id  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_home_product_%s (date  datetime,uin  int unsigned,petid  varchar(255),house_id  int,goods_id  varchar(255),price  int,vip  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_gift_%s (date  datetime,uin  int unsigned,petid  varchar(255),act_id  int,opr_id  int,giftid  int,index(act_id,opr_id,giftid),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_soc_game_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,vip  int,game_id  int,scene_id  int,beg_end  varchar(255),yb_val  varchar(255),qb_val  varchar(255),starve  varchar(255),clean  varchar(255),strong  varchar(255),iq  varchar(255),charm  varchar(255),bill_id  varchar(255),level  varchar(255),plat_id  varchar(255),index(game_id,scene_id),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_soc_game_goods_%s (date  datetime,uin  int unsigned,petid  varchar(255),game_id  int,bill_id  varchar(255),goods_id  varchar(255),goods_num  varchar(255),plat_id  varchar(255),index(game_id, goods_id),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_soc_oss_cm_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,game_id  int,receive_uin  int unsigned,game_name  varchar(255),detail_desc  varchar(255),index(game_id),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_soc_game_err_%s (date  datetime,uin  int unsigned,petid  varchar(255),game_id  int,err_code  varchar(255),detail_info  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_soc_goods_%s (date  datetime,uin  int unsigned,petid  varchar(255),game_id  int,goods_id  varchar(255),goods_num  varchar(255),vip  varchar(255),scene_id  varchar(255),beg_end  varchar(255),index(game_id),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tbReg_%s (date  datetime,ip  varchar(255),ret  int,way  int,uin  int unsigned,sex  int,petid  varchar(255),father_id  varchar(255),mother_id  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tbGrowthValueChg_%s (date  datetime,chg_id  int,uin  int unsigned,old_value  int,last_chg_time  int,level  int,chg_value  int,new_value  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_trip_gift_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),goods_id  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_trip_start_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),set_type  varchar(255),trip_mins  varchar(255),map_id  varchar(255),city_id  varchar(255),city_name  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_trip_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),set_type  varchar(255),trip_mins  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_adopt_svrinfo_%s (date  datetime,uin  int unsigned,pet_num  int,pet_type_num  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_kernel_version_%s (date  datetime,uin  int unsigned,typeid  varchar(255),kernel_version  varchar(255),local_version  varchar(50),index(typeid, local_version),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_fish_%s (date  datetime,uin  int unsigned,petid  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tbLogout_%s (date  datetime,uin  int unsigned,sex  int,petid  varchar(255),ret  int,last_time  varchar(255),ip  varchar(255),login_type  int,online_time  int,total_time  int,effect  int,level  int,vip  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_pet_discard_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),ret  int,health  int,sick_id  int,level  int,name  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_feed_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,feed_channel  varchar(255),goods_id  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_attrib_change_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),procedure_id  varchar(255),growth  varchar(255),starve  varchar(255),clean  varchar(255),mood  varchar(255),strong  varchar(255),brain  varchar(255),charm  varchar(255),love  varchar(255),trip  varchar(255),goods_id  varchar(255),effective_flag  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_avatar_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),ret  int,Ava1  varchar(255),Ava2  varchar(255),Ava3  varchar(255),Ava4  varchar(255),Ava5  varchar(255),Ava6  varchar(255),Ava7  varchar(255),Ava8  varchar(255),Ava9  varchar(255),Ava10  varchar(255),Ava11  varchar(255),Ava12  varchar(255),Ava13  varchar(255),Ava14  varchar(255),Ava15  varchar(255),Ava16  varchar(255),Ava17  varchar(255),Ava18  varchar(255),Ava19  varchar(255),Ava20  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_yb_change_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),procedure_id  varchar(255),num  int,remain_num  int,index(procedure_id, num),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_goods_change_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),procedure_id  varchar(255),goods_id  varchar(255),goods_num  int,goods_exp  int,goods_num_remain  int,goods_exp_remain  int,index(goods_id),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_home_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,guest_uin  int unsigned,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_time_out_%s (date  datetime,uin  int unsigned,petid  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_work_start_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),work_id  int,work_mins  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_work_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),work_id  varchar(50),work_mins  varchar(50),work_name  varchar(255),strong  varchar(255),brain  varchar(255),charm  varchar(255),index(work_id, uin, work_mins),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_sick_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),health  varchar(255),sick_id  varchar(255),starve  varchar(255),clean  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_study_start_%s (date  datetime,uin  int unsigned,petid  varchar(255),study_id  varchar(255),study_mins  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_study_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),study_id  varchar(255),remain_lessons  varchar(255),strong  varchar(255),brain  varchar(255),charm  varchar(255),lesson  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_pet_state_start_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),stat_id  varchar(255),stat_num  varchar(255),stat_mins  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_pet_state_end_%s (date  datetime,uin  int unsigned,species  int,petid  varchar(255),stat_id  varchar(255),health_flag  varchar(255),yb_change  varchar(255),index(stat_id, uin),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_zone_switch_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  int,zone_id  int,area_id  int,pre_scene_id  int,scene_id  int,index(scene_id,uin),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_darling_state_%s (date  datetime,uin  int unsigned,sex  int,petid  varchar(255),success  varchar(255),wrong  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_autofeed_%s (date  datetime,uin  int unsigned,sex  int,petid  varchar(255),success  varchar(255),goods_id  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_hungryfeed_%s (date  datetime,uin  int unsigned,sex  int,petid  varchar(255),success  varchar(255),goods_id  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_chess_use_%s (date  datetime,log_type  varchar(255),ret  varchar(255),detail_info  varchar(255),table_type  varchar(255),1st_uin  int unsigned,2nd_uin  int unsigned,uin  int unsigned,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_chess_buy_%s (date  datetime,log_type  varchar(255),uin  int unsigned,petid  varchar(255),retCode  varchar(255),detail_info  varchar(255),qb  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_chess_start_%s (date  datetime,log_type  varchar(255),uin  int unsigned,retCode  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_chess_game_%s (date  datetime,log_type  varchar(255),uin  int unsigned,petid  varchar(255),retCode  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_chess_score_%s (date  datetime,log_type  varchar(255),uin  int unsigned,petid  varchar(255),ret  varchar(255),detail_info  varchar(255),score_change  varchar(255),score  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_zone_%s (date  datetime,uin  int unsigned,petid  varchar(255),ret  varchar(255),zone_id  varchar(255),area_id  varchar(255),scene_id  int,online_sum  varchar(255),level  int,vip  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
CREATE TABLE if not exists petLog_%d.tb_community_%s (date datetime, uin int unsigned, status_id int, loadingtimes int, index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_idip_%s (date  datetime,uin  int unsigned,zoneid  int unsigned,item_id  varchar(255),item_count  int, idip_cmd  int unsigned,seq  varchar(255),src_id  int unsigned,act_id  int unsigned,ip  varchar(255),oper  varchar(255),req_str  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.bblog_pay_%s (date  datetime,uin  int unsigned,platform  int,vipflag  int,viplevel  int,vipyear  int,sid  int,qb  int,petjuan  int,num  int,good_type  int,good_p1  int,good_p2  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.bblog_task_%s (date  datetime,uin  int unsigned,platform  int,vipflag  int,viplevel  int,vipyear  int,task_type  int,task_id  int,pre_status  int,status  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.bblog_change_%s (date  datetime,uin  int unsigned,platform  int,vipflag  int,viplevel  int,vipyear  int,chg_type  int,chg_id  int,chg_p1  int,chg_p2  int,chg_reason  varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_send_msg_%s (date  datetime,uin  int unsigned,type  int,msgid  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_interact_%s (date  datetime,uin  int unsigned,type  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_stealth_time_%s (date  datetime,uin  int unsigned,start  int unsigned,time  int unsigned,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_start_status_%s (date  datetime,uin  int unsigned,ret  int,ip  varchar(255),status  int,last_status  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.t_stop_status_%s (date  datetime,uin  int unsigned,ret  int,last_status  int,over_time  int,index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
create table if not exists petLog_%d.tb_petquan_change_%s (date  datetime,uin  int unsigned,ip  varchar(255),module_id  int unsigned,change_num  int,remain_num  int unsigned, desc_info varchar(255),index(uin)) ENGINE=InnoDB DEFAULT CHARSET=gbk;
