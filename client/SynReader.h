#ifndef SYN_READER_H
#define SYN_READER_H

#ifdef __cplusplus
extern "C" {
#endif



#pragma pack(1)
//普通居民身份证或港澳台居民居住证
typedef struct IDCardData { 
    unsigned short Name[15];         // offset=0  
    unsigned short Sex;              // 30 
    unsigned short Nation[2];        // 32 
    unsigned short Birthday[8];      // 36 
    unsigned short Address[35];      // 52
    unsigned short IDCardNo[18];     // 122 
    unsigned short GrantDept[15];    // 158 
    unsigned short UserLifeBegin[8]; // 188 
    unsigned short UserLifeEnd[8];   // 204 
 	unsigned short PassID[9];   	//220  
 	unsigned short IssuesTimes[2];	//238 
	unsigned short reserved1[3];	//242
	unsigned short CardType;		//248
	unsigned short reserved2[3];//
}St_IDCardData, *PSt_IDCardData;

//外国人永久居留证
typedef struct ForeignerCardData { 
    unsigned short EngName[60];		// offset=0  
    unsigned short Sex;				// 120 
 	unsigned short IDCardNo[15];	// 122 
    unsigned short Nation[3]; 		// 152 
 	unsigned short Name[15];  		//158
    unsigned short UserLifeBegin[8];// 188
    unsigned short UserLifeEnd[8];  // 204 
    unsigned short Birthday[8];     // 220 
	unsigned short CertVol[2];     	//	236
    unsigned short GrantDept[4];    // 240
	unsigned short CardType;		//248
	unsigned short reserved2[3];
} St_ForeignerCardData, *PSt_ForeignerCardData;



typedef struct IDCardDataUTF8 {
	char CardType[10]; 		//I为外国人居住证，J 为港澳台居住证，空格(0x20)为普通身份证
    char Name[40];         	//姓名 
	char EngName[130];   	//英文名(外国人居住证)
    char Sex[10];           //性别
    char Nation[100];  		//民族或国籍(外国人居住证)     
    char Birthday[18];     	//出生日期
    char Address[80];       //住址
    char IDCardNo[40];      //身份证号或外国人居住证号(外国人居住证)
    char GrantDept[40];     //发证机关
    char UserLifeBegin[30]; //有效开始日期
    char UserLifeEnd[30];   //有效截止日期
	char PassID[30];		//通行证号码(港澳台)
	char IssuesTimes[10];	//签发次数(港澳台)
	char CertVol[10];     	//证件版本号(外国人居住证)
    char wlt[1024];   		//照片数据
    int isSavePhotoOK;		//照片是否解码保存  0=no  1=yes
	char fp[1024];			//指纹数据
	int isFpRead;			//是否读取了证内指纹	 0=no 1=yes
} St_IDCardDataUTF8, *PSt_IDCardDataUTF8;
#pragma pack()
//寻卡 pucIIN 4个字节
int SAM_FindCard(unsigned char * pucIIN);	
//选卡 pucSN 8个字节
int SAM_SelectCard(unsigned char * pucSN);
//读卡 pucBaseMsg 文字信息 256字节 pucPhoto照片原始数据1024字节
int SAM_ReadBaseMsg(unsigned char * pucBaseMsg,unsigned char * pucPhoto);
//读卡和证内指纹 pucBaseMsg 文字信息 256字节 pucPhoto照片原始数据1024字节 pucFpData指纹数据1024字节
int SAM_ReadBaseFpMsg(unsigned char * pucBaseMsg,unsigned char * pucPhoto,unsigned char * pucFpData);

//读身份证文字和照片信息
int getIDcard(St_IDCardDataUTF8 *pIDCardDataUTF8);
//读身份证文字、照片和指纹信息
int getIDcardWithFp(St_IDCardDataUTF8 *pIDCardDataUTF8);
//将pucBaseMsg转换为St_IDCardDataUTF8
int baseMsg2IDCardDataUTF8(unsigned char * pucBaseMsg,St_IDCardDataUTF8 *pIDCardDataUTF8);



typedef enum KeyType{KeyA,KeyB} KeyType;
int m1FindCardWithNo(unsigned char* cardNo);
int m1FindCard();
int m1DownloadKey(KeyType type,unsigned char block,const unsigned char *key);
int m1CheckKey(KeyType type,unsigned char block);
int m1ReadBlock(unsigned char block,unsigned char* blockData);
int m1WriteBlock(unsigned char block,const unsigned char* blockData);
int m1Pend();
void printArray(const unsigned char *data,int len);

int Syn_USBHIDM1Reset(unsigned char * cardInfo, int * infoLen);
int Syn_USBHIDM1AuthenKey(KeyType type, unsigned char BlockNo, unsigned char * pKey);
int Syn_USBHIDM1ReadBlock(unsigned char BlockNo, unsigned char * pBlock);
int Syn_USBHIDM1WriteBlock(unsigned char BlockNo, unsigned char * pBlock);
int Syn_USBHIDM1Halt();



typedef enum CommType{Serial,USB,NONE}CommType;
int OpenUsbComm();
int OpenSerialComm(const char*);
int CloseComm();


//解码照片函数，wltBuffer原始照片1024字节wlt数据，bmpPath 保存照片位置
int saveWlt2Bmp( char* wltBuffer,const char* bmpPath);
int saveWlt2BmpUseFork( char* wltBuffer,const char* bmpPath);

//获取动态库版本
const char* getLibVersion();
int getLibVersionInt();



#ifdef __cplusplus
}
#endif

#endif
