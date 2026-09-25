#include <iostream>
#include <string_view>

/* C++ block comment
   with a second line. */
int main() {
    constexpr std::string_view text = R"tag(raw "value")tag";
    std::cout << text << '\n';
}
