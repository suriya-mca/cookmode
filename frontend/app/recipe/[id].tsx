import { View, Text } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { theme } from "@/lib/theme";

export default function RecipeDetail() {
  const { id } = useLocalSearchParams<{ id: string }>();
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        padding: theme.spacing.lg,
      }}
    >
      <Text style={{ fontSize: theme.text.title, color: theme.colors.charcoal }}>
        Recipe {id}
      </Text>
      <Text style={{ marginTop: 8, color: theme.colors.charcoalMuted }}>
        GET /recipes/{id} → hero video, ingredients (servings stepper), steps with anchor
      </Text>
    </View>
  );
}
