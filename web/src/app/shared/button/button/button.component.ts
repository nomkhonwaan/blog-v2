import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { IconDefinition as SolidIconDefinition } from '@fortawesome/fontawesome-free-solid';
import { IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { IconDefinition } from '@fortawesome/pro-light-svg-icons';

@Component({
  selector: 'app-button',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.scss'],
})
export class ButtonComponent {

  /**
   * A FontAwesome icon to-be displayed with the button
   */
  @Input()
  icon: IconDefinition | SolidIconDefinition | BrandIconDefinition;

  /**
   * An HTML aria-label property
   */
  @Input()
  ariaLabel: string;

  hasIcon(): boolean {
    return !!this.icon;
  }

}
