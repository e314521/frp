#include "SynReader.h"
#include <stdio.h>
#include <string.h>
int ZKID_Init(){
  return 0;
}
int ZKID_OpenPort(int Port){
  return -2;
}
int ZKID_Free(){
  return 0;
}
int ZKID_GetSAMStatus(int Port, int iIfOpen){
  return 0;
}
int ZKID_GetSAMIDToStr(int Port, char*pcSAMID, int iIfOpen){
  memcpy(pcSAMID,"123",3);
  return 0;
}
int ZKID_StartFindIDCard(int Port, char*pcSAMID, int iIfOpen){
  return -1;
}
