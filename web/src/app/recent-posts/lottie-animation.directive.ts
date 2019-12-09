import { Directive, OnInit, Input, ElementRef, NgZone } from '@angular/core';
import Lottie from 'lottie-web';

@Directive({ selector: '[appLottieAnimation]' })
export class LottieAnimationDirective implements OnInit {

  @Input()
  data: any;

  constructor(private elementRef: ElementRef, private ngZone: NgZone) { }

  ngOnInit(): void {
    this.ngZone.runOutsideAngular((): void => {
      Lottie.loadAnimation({
        container: this.elementRef.nativeElement,
        renderer: 'svg',
        loop: true,
        autoplay: true,
        animationData: this.data,
      });
    });
  }

}
