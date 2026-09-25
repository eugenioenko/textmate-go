#import <Foundation/Foundation.h>
#include <vector>

int main(void) {
    std::vector<int> values{1, 2, 3};
    NSLog(@"count: %zu", values.size());
}
