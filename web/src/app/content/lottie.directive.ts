import { Directive, ElementRef, Input, NgZone, OnInit } from '@angular/core';
import Lottie from 'lottie-web';

@Directive({ selector: '[appLottie]' })
export class LottieDirective implements OnInit {

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
