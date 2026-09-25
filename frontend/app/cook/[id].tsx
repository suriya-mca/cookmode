import { useEffect, useRef, useState } from "react";
import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import { VideoView, useVideoPlayer } from "expo-video";
import { useKeepAwake } from "expo-keep-awake";
import * as Speech from "expo-speech";
import BottomSheet, { BottomSheetView } from "@gorhom/bottom-sheet";
import { GestureHandlerRootView } from "react-native-gesture-handler";
import { useRecipe } from "@/lib/recipes";
import { ApiError } from "@/lib/api";
import { theme } from "@/lib/theme";

export default function CookMode() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { data: r, isLoading, error } = useRecipe(id!);
  const [stepIdx, setStepIdx] = useState(0);
  const [timerSec, setTimerSec] = useState<number | null>(null);
  const sheetRef = useRef<BottomSheet>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useKeepAwake();

  const hlsUrl = r?.video_hls_url || undefined;
  const player = useVideoPlayer(hlsUrl ?? null, (p: any) => {
    p.loop = false;
  });

  useEffect(() => {
    if (!r || !hlsUrl) return;
    try {
      player.replaceAsync({ uri: hlsUrl } as any);
    } catch {}
  }, [hlsUrl]);

  useEffect(() => {
    if (timerSec === null) return;
    if (timerSec <= 0) {
      Speech.speak("Timer done");
      setTimerSec(null);
      if (intervalRef.current) clearInterval(intervalRef.current);
      return;
    }
    intervalRef.current = setInterval(() => {
      setTimerSec((s) => (s !== null ? s - 1 : null));
    }, 1000);
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [timerSec]);

  const steps = r?.steps ?? [];
  const step = steps[stepIdx];
  const anchor = step?.anchor_seconds ?? 0;

  const seekTo = (sec: number) => {
    try {
      (player as any).currentTime = sec;
      player.play();
    } catch {}
  };

  const next = () => setStepIdx((i) => Math.min(i + 1, steps.length - 1));
  const prev = () => setStepIdx((i) => Math.max(i - 1, 0));
  const repeat = () => seekTo(anchor);
  const setTimer = (sec: number) => setTimerSec(sec);

  if (isLoading) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.charcoal, alignItems: "center", justifyContent: "center" }}>
        <ActivityIndicator color={theme.colors.cream} />
      </View>
    );
  }
  // A stale token turns this public read into 401; don't claim the recipe is
  // missing — offer a re-sign-in instead.
  if (error instanceof ApiError && error.status === 401) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.charcoal, alignItems: "center", justifyContent: "center", padding: 24 }}>
        <Text style={{ color: theme.colors.cream, fontWeight: "600" }}>Session expired</Text>
        <Text style={{ color: theme.colors.cream, opacity: 0.7, fontSize: 13, marginTop: 8, textAlign: "center" }}>
          Sign in again to start cooking.
        </Text>
        <Pressable
          onPress={() => router.replace("/(auth)/login")}
          style={{ marginTop: 16, paddingHorizontal: 18, paddingVertical: 12, backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm }}
        >
          <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Sign in</Text>
        </Pressable>
      </View>
    );
  }
  if (!r) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.charcoal, alignItems: "center", justifyContent: "center" }}>
        <Text style={{ color: theme.colors.cream }}>Recipe not found</Text>
        <Pressable onPress={() => router.back()} style={{ marginTop: 12, padding: 12, backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm }}>
          <Text style={{ color: theme.colors.card }}>Go back</Text>
        </Pressable>
      </View>
    );
  }

  return (
    <GestureHandlerRootView style={{ flex: 1, backgroundColor: theme.colors.charcoal }}>
      <View style={{ flex: 1 }}>
        <View style={{ height: "60%", backgroundColor: "#000" }}>
          {hlsUrl ? (
            <VideoView player={player} style={{ flex: 1 }} contentFit="contain" nativeControls />
          ) : (
            <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
              <Text style={{ color: theme.colors.cream, opacity: 0.7 }}>No video — steps only</Text>
            </View>
          )}
          <View style={{ position: "absolute", top: 48, left: 16, right: 16, flexDirection: "row", justifyContent: "space-between", alignItems: "center" }}>
            <Pressable onPress={() => router.back()} style={{ backgroundColor: "rgba(0,0,0,0.5)", paddingHorizontal: 12, paddingVertical: 6, borderRadius: 20 }}>
              <Text style={{ color: theme.colors.card }}>✕ Close</Text>
            </Pressable>
            <Text style={{ color: theme.colors.card, backgroundColor: "rgba(0,0,0,0.5)", paddingHorizontal: 10, paddingVertical: 6, borderRadius: 20, overflow: "hidden" }}>
              Step {stepIdx + 1}/{steps.length || 1}
            </Text>
            {timerSec !== null ? (
              <Text style={{ color: theme.colors.card, backgroundColor: theme.colors.sage, paddingHorizontal: 10, paddingVertical: 6, borderRadius: 20, overflow: "hidden" }}>
                {Math.floor(timerSec / 60)}:{String(timerSec % 60).padStart(2, "0")}
              </Text>
            ) : (
              <View style={{ width: 60 }} />
            )}
          </View>
        </View>

        <BottomSheet
          ref={sheetRef}
          index={1}
          snapPoints={["35%", "65%"]}
          backgroundStyle={{ backgroundColor: theme.colors.card, borderTopLeftRadius: theme.radius.lg, borderTopRightRadius: theme.radius.lg }}
          handleIndicatorStyle={{ backgroundColor: theme.colors.border }}
        >
          <BottomSheetView style={{ flex: 1, padding: 16, gap: 12 }}>
            <Text style={{ fontSize: 12, color: theme.colors.terracotta, fontWeight: "700", letterSpacing: 1 }}>STEP {String(stepIdx + 1).padStart(2, "0")}</Text>
            <Text style={{ fontSize: 20, color: theme.colors.charcoal, fontWeight: "700", lineHeight: 26 }}>
              {step?.text ?? "No steps yet — add steps in the recipe."}
            </Text>
            <Text style={{ color: theme.colors.sage, fontSize: 13 }}>▶ {anchor}s • {step?.duration_hint ? `~ ${step.duration_hint} min` : "tap Repeat to seek"}</Text>

            <View style={{ flexDirection: "row", gap: 8, marginTop: 8 }}>
              <Pressable onPress={prev} disabled={stepIdx === 0} style={{ flex: 1, backgroundColor: theme.colors.creamSoft, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 14, alignItems: "center", opacity: stepIdx === 0 ? 0.5 : 1 }}>
                <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>Previous</Text>
              </Pressable>
              <Pressable onPress={next} disabled={stepIdx >= steps.length - 1} style={{ flex: 1, backgroundColor: theme.colors.charcoal, borderRadius: theme.radius.sm, padding: 14, alignItems: "center", opacity: stepIdx >= steps.length - 1 ? 0.5 : 1 }}>
                <Text style={{ fontWeight: "600", color: theme.colors.card }}>Next Step →</Text>
              </Pressable>
            </View>

            <View style={{ flexDirection: "row", gap: 8 }}>
              <Pressable onPress={repeat} style={{ flex: 1, backgroundColor: theme.colors.cream, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.pill, padding: 12, alignItems: "center" }}>
                <Text style={{ color: theme.colors.charcoal, fontWeight: "600" }}>Repeat ▶ {anchor}s</Text>
              </Pressable>
              <Pressable onPress={() => setTimer(240)} style={{ flex: 1, backgroundColor: theme.colors.sage, borderRadius: theme.radius.pill, padding: 12, alignItems: "center" }}>
                <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Set Timer 4:00</Text>
              </Pressable>
            </View>

            <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12, textAlign: "center", marginTop: 4 }}>
              Voice: “next step” • “repeat” • “set timer” (stub — expo-speech ready)
            </Text>
          </BottomSheetView>
        </BottomSheet>
      </View>
    </GestureHandlerRootView>
  );
}
