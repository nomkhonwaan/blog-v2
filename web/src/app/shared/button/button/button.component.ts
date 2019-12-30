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
   * An HTML aria-label property
   */
  @Input()
  ariaLabel: string;

  /**
   * A FontAwesome icon to-be displayed with the button
   */
  @Input()
  icon: IconDefinition | SolidIconDefinition | BrandIconDefinition;

  /**
   * Use to tell FontAwesome component to spin the icon or not
   */
  @Input()
  spin: boolean;

  /**
   * Use to display button color in "danger" mode
   */
  @Input()
  danger: boolean;

  hasIcon(): boolean {
    return !!this.icon;
  }

  isDanger(): boolean {
    return !!this.danger;
  }

}
