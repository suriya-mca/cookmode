import { View, Text } from "react-native";
import { theme } from "@/lib/theme";

export default function AddScreen() {
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Text style={{ color: theme.colors.charcoal, fontSize: theme.text.subtitle }}>
        Add Recipe
      </Text>
      <Text style={{ color: theme.colors.charcoalMuted, marginTop: 8 }}>
        POST /recipes + upload via POST /recipes/upload-url
      </Text>
    </View>
  );
}
