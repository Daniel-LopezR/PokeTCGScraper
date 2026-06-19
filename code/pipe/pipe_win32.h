#define BUFSIZE 4096
#define CMD_BUFSIZE 512

void CreatePipes(void);
void CreateChildProcess(void);
void ReadFromPipe(void);
void WriteToCmdPipe(char[BUFSIZE]);
