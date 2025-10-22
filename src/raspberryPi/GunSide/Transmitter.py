import lgpio
import time


def shootWithInfo():
    pass

def shoot():
    GPIO_PIN = 18
    FREQUENCY = 38000
    DUTY_CYCLE = 50

    PULSE_WIDTH_US = int(1_000_000 / FREQUENCY / 2)

    if PULSE_WIDTH_US == 0:
        PULSE_WIDTH_US == 1
    OFF_WIDTH_US = PULSE_WIDTH_US

    TRANSMIT_DURATION = 0.5

    try:
        h = lgpio.gpiochip_open(0)

        lgpio.gpio_claim_output(h, GPIO_PIN)

        p1 = lgpio.pulse(1 << GPIO_PIN, 0, PULSE_WIDTH_US)
        p2 = lgpio.pulse(0, 1 << GPIO_PIN, OFF_WIDTH_US)

        wave = [p1, p2]

        lgpio.wave_clear(h)
        lgpio.wave_add_generic(h, wave)
        wave_id = lgpio.wave_create(h)

        if wave_id >= 0:

            lgpio.wave_tx_send(h, wave_id, lgpio.WAVE_MODE_REPEAT)

            time.sleep(TRANSMIT_DURATION)

            lgpio.wave_tx_stop(h)
            lgpio.wave_delete(h, wave_id)
        else:
            print(f"error with creating wave: {wave_id}")
    except Exception as e:
        print(f"an exception has occured: \n {e}")
    finally:
        if 'h' in locals() and h >= 0:
            lgpio.gpio_free(h, GPIO_PIN)
            lgpio.gpiochip_close(h)
    