#ifdef _WIN32
#include "pipe_win32.h"
#include "windows.h"
#include "stdio.h"
#include "stdbool.h"
#include "stdint.h"
#include <minwindef.h>
#include <stdint.h>
#include <string.h>

HANDLE global_cmdHandle_Rd = NULL;
HANDLE global_cmdHandle_Wr = NULL;
HANDLE global_dataHandle_Rd = NULL;
HANDLE global_dataHandle_Wr = NULL;

void CreatePipes()
{
    SECURITY_ATTRIBUTES sec_attr;

    sec_attr.nLength = sizeof(SECURITY_ATTRIBUTES); 
    sec_attr.bInheritHandle = TRUE; 
    sec_attr.lpSecurityDescriptor = NULL;
    
    if (!CreatePipe(&global_cmdHandle_Rd, &global_cmdHandle_Wr, &sec_attr, 0)){
    	DWORD err = GetLastError();
    	printf("Couldn't create pipe for commands: %lu\n", err);
    }
    printf("Commands Pipe created\n");
    
    if (!SetHandleInformation(global_cmdHandle_Wr, HANDLE_FLAG_INHERIT, 0)) {
	    DWORD err = GetLastError();
	    printf("Coudldn't remove HANDLE_FLAG_INHERIT from global_cmdHandle_Wr:  %lu\n", err);
    }
    printf("Removed HANDLE_FLAG_INHERIT from global_cmdHandle_Wr\n");
    
    if (!CreatePipe(&global_dataHandle_Rd, &global_dataHandle_Wr, &sec_attr, 0)){
    	DWORD err = GetLastError();
    	printf("Couldn't create pipe for data: %lu\n", err);
    }
    printf("Data Pipe created\n");
    
    if (!SetHandleInformation(global_dataHandle_Rd, HANDLE_FLAG_INHERIT, 0)) {
    	DWORD err = GetLastError();
    	printf("Coudldn't remove HANDLE_FLAG_INHERIT from global_dataHandle_Rd:  %lu\n", err);
    }
    printf("Removed HANDLE_FLAG_INHERIT from global_dataHandle_Rd\n");
};
void CreateChildProcess()
// Create a child process that uses the previously created pipes for STDIN and STDOUT.
{
    TCHAR szCmdline[]=TEXT("child.exe");
    PROCESS_INFORMATION piProcInfo; 
    STARTUPINFO siStartInfo;
    BOOL bSuccess = FALSE; 
    
    // Set up members of the PROCESS_INFORMATION structure. 
     
    ZeroMemory( &piProcInfo, sizeof(PROCESS_INFORMATION) );
    
    // Set up members of the STARTUPINFO structure. 
    // This structure specifies the STDIN and STDOUT handles for redirection.
    
    // Fills a block of memory with zeros.
    ZeroMemory( &siStartInfo, sizeof(STARTUPINFO) );
    siStartInfo.cb = sizeof(STARTUPINFO);
    // Maybe create a separate handle for erros to avoid confusions
    siStartInfo.hStdError = global_dataHandle_Wr;
    siStartInfo.hStdOutput = global_dataHandle_Wr;
    siStartInfo.hStdInput = global_cmdHandle_Rd;
    siStartInfo.dwFlags |= STARTF_USESTDHANDLES;
     
    // Create the child process. 
    bSuccess = CreateProcess(NULL, 
        szCmdline,     // command line 
        NULL,          // process security attributes 
        NULL,          // primary thread security attributes 
        TRUE,          // handles are inherited 
        0,             // creation flags 
        NULL,          // use parent's environment 
        NULL,          // use parent's current directory 
        &siStartInfo,  // STARTUPINFO pointer 
        &piProcInfo);  // receives PROCESS_INFORMATION 
   
    // If an error occurs, exit the application. 
    if ( ! bSuccess ) {
	    DWORD err = GetLastError();
	    printf("Coudln't create child process: %lu\n", err);
    } else {
        // Close handles to the child process and its primary thread.
        // Some applications might keep these handles to monitor the status
        // of the child process, for example. 
    
        CloseHandle(piProcInfo.hProcess);
        CloseHandle(piProcInfo.hThread);
         
        // Close handles to the stdin and stdout pipes no longer needed by the child process.
        // If they are not explicitly closed, there is no way to recognize that the child process has ended.
        
        CloseHandle(global_dataHandle_Wr);
	    CloseHandle(global_cmdHandle_Rd);
    }
}

void ReadFromPipe(void) 
// Read output from the child process's pipe for STDOUT
// and write to the parent process's pipe for STDOUT. 
// Stop when there is no more data. 
{
    DWORD dwRead, dwWritten;
    LPDWORD dwPeeked;
    CHAR chBuf[BUFSIZE];
    BOOL bSuccess = FALSE;
    HANDLE hParentStdOut = GetStdHandle(STD_OUTPUT_HANDLE);
    
    printf("Reading from pipe...\n");
    // The last 2 params will be important to read to manage cases where more messages are left to read or the message
    // itself is greater than the buffer, which should not be posible but jsut in case
    //[out, optional] LPDWORD lpTotalBytesAvail,
    //[out, optional] LPDWORD lpBytesLeftThisMessage
    bSuccess = PeekNamedPipe(global_dataHandle_Rd, chBuf, BUFSIZE, dwPeeked, NULL, NULL);
    
    if(!bSuccess){
        DWORD err = GetLastError();
        printf("Coudldn't peek data from global_dataHandle_Rd:  %lu\n", err);
        return;
    }

    if (dwPeeked != 0){
        
	    bSuccess = ReadFile(global_dataHandle_Rd, chBuf, BUFSIZE, &dwRead, NULL);
        if(!bSuccess){
	        DWORD err = GetLastError();
	        printf("Coudldn't remove HANDLE_FLAG_INHERIT from global_cmdHandle_Wr:  %lu\n", err);
            return;
        }  
   
        if (dwRead != 0) {
            //chBuf[BUFSIZE-1] = '\n';
            bSuccess = WriteFile(hParentStdOut, chBuf, dwRead, &dwWritten, NULL);
            if (!bSuccess){
	            DWORD err = GetLastError();
	            printf("Coudldn't remove HANDLE_FLAG_INHERIT from global_cmdHandle_Wr:  %lu\n", err);
                return;
            }
            //printf("\n");
        }
    }
}

void WriteToCmdPipe(char cmd[CMD_BUFSIZE])
{
    char chBuf[CMD_BUFSIZE] = {0};
    // The above acomplishes the same, but not sure about the implications of each method
    // ZeroMemory(chBuf, CMD_BUFSIZE);
    strcpy_s(chBuf, CMD_BUFSIZE, cmd);
    DWORD dwWritten; 
    bool bSuccess = FALSE;
    // Use the debugger to look into the memory layout
    chBuf[CMD_BUFSIZE-1] = '\n';
    printf("Writing to cmd pipe...\n");
    
    bSuccess = WriteFile(global_cmdHandle_Wr, chBuf, CMD_BUFSIZE, &dwWritten, NULL); 
    if (dwWritten != 0) {
       printf("Bytes written from cmd buffer -> %lu\n", dwWritten); 
        if (!bSuccess) {
	        DWORD err = GetLastError();
	        printf("Coudldn't write to cmd pipe:  %lu\n", err);
            return;
        }
    }
} 
#endif
