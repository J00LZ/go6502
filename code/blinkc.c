
void loop();


void reset() {
    while (1) {
        loop();
    }
}


int a = 0;
void loop() {
    a += 1;
}